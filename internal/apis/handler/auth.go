package handler

import (
	"context"
	"errors"

	"github.com/QuizWars-Ecosystem/users-service/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/status"

	"github.com/QuizWars-Ecosystem/go-common/pkg/abstractions"
	apperrors "github.com/QuizWars-Ecosystem/go-common/pkg/error"
	"github.com/QuizWars-Ecosystem/go-common/pkg/jwt"
	"github.com/QuizWars-Ecosystem/go-common/pkg/uuidx"
	"go.uber.org/zap"

	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/auth"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) Register(ctx context.Context, request *userspb.RegisterRequest) (*userspb.RegisterResponse, error) {
	var err error
	req, err := abstractions.MakeRequest[auth.ProfileWithCredentials](request)
	if err != nil {
		return nil, err
	}

	defer metrics.UserCreatedTotalCounter.WithLabelValues("server", status.Code(err).String()).Inc()
	timer := prometheus.NewTimer(metrics.UsersCreationDurationHistogram.WithLabelValues("server"))
	defer timer.ObserveDuration()
	defer takeErrorMetric("server", err)

	res, err := h.service.Register(ctx, req)
	if err != nil {
		return nil, err
	}

	result, err := abstractions.MakeResponse(res)
	if err != nil {
		return nil, err
	}

	token, err := h.jwt.GenerateToken(res.User.ID.String(), string(jwt.User))
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	h.logger.Debug("new user registered", zap.String("id", result.Id))

	return &userspb.RegisterResponse{
		Token:   "Bearer " + token,
		Profile: result,
	}, nil
}

func (h *Handler) Login(ctx context.Context, request *userspb.LoginRequest) (*userspb.LoginResponse, error) {
	var credits *auth.ProfileWithCredentials
	var err error

	defer metrics.UsersLoginTotalCounter.WithLabelValues("server", status.Code(err).String()).Inc()
	timer := prometheus.NewTimer(metrics.UsersLoginDurationHistogram.WithLabelValues("server"))
	defer timer.ObserveDuration()

	switch request.Identifier.(type) {
	case *userspb.LoginRequest_Username:
		credits, err = h.service.LoginByUsername(ctx, request.GetUsername(), request.GetPassword())
	case *userspb.LoginRequest_Email:
		credits, err = h.service.LoginByEmail(ctx, request.GetEmail(), request.GetPassword())
	}

	if err != nil {
		return nil, err
	}

	result, err := abstractions.MakeResponse(credits.Profile)
	if err != nil {
		return nil, err
	}

	token, err := h.jwt.GenerateToken(credits.Profile.User.ID.String(), credits.Role)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	return &userspb.LoginResponse{
		Token:   "Bearer " + token,
		Profile: result,
	}, nil
}

func (h *Handler) Logout(ctx context.Context, request *userspb.LogoutRequest) (*emptypb.Empty, error) {
	userID, err := uuidx.Parse(request.GetUserId())
	if err != nil {
		return nil, err
	}

	err = h.service.Logout(ctx, userID)
	if err != nil {
		return nil, err
	}

	metrics.UserLogoutTotalCounter.Inc()

	return Empty, nil
}

func (h *Handler) OAuthLogin(_ context.Context, _ *userspb.OAuthLoginRequest) (*userspb.OAuthLoginResponse, error) {
	return nil, apperrors.Internal(errors.New("not implemented"))
}

func (h *Handler) LinkOAuthProvider(_ context.Context, _ *userspb.LinkOAuthProviderRequest) (*emptypb.Empty, error) {
	return Empty, apperrors.Internal(errors.New("not implemented"))
}

func takeErrorMetric(method string, err error) {
	if err != nil {
		metrics.UsersCreationErrorsCounter.WithLabelValues(method, status.Code(err).String()).Inc()
	}
}
