package handler

import (
	"context"
	"errors"

	"github.com/QuizWars-Ecosystem/go-common/pkg/abstractions"
	apperrors "github.com/QuizWars-Ecosystem/go-common/pkg/error"
	"github.com/QuizWars-Ecosystem/go-common/pkg/jwt"
	"go.uber.org/zap"

	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/auth"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) Register(ctx context.Context, request *userspb.RegisterRequest) (*userspb.RegisterResponse, error) {
	req, err := abstractions.MakeRequest[auth.ProfileWithCredentials](request)
	if err != nil {
		return nil, err
	}

	res, err := h.service.Register(ctx, req)
	if err != nil {
		return nil, err
	}

	result, err := abstractions.MakeResponse(res)
	if err != nil {
		return nil, err
	}

	token, err := h.jwt.GenerateToken(res.User.ID, string(jwt.User))
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

	token, err := h.jwt.GenerateToken(credits.Profile.User.ID, credits.Role)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	return &userspb.LoginResponse{
		Token:   "Bearer " + token,
		Profile: result,
	}, nil
}

func (h *Handler) Logout(ctx context.Context, request *userspb.LogoutRequest) (*emptypb.Empty, error) {
	err := h.service.Logout(ctx, request.GetUserId())
	if err != nil {
		return nil, err
	}

	return Empty, nil
}

func (h *Handler) OAuthLogin(_ context.Context, _ *userspb.OAuthLoginRequest) (*userspb.OAuthLoginResponse, error) {
	return nil, apperrors.Internal(errors.New("not implemented"))
}

func (h *Handler) LinkOAuthProvider(_ context.Context, _ *userspb.LinkOAuthProviderRequest) (*emptypb.Empty, error) {
	return Empty, apperrors.Internal(errors.New("not implemented"))
}
