package handler

import (
	"context"
	apperrors "github.com/QuizWars-Ecosystem/go-common/pkg/error"
	"github.com/QuizWars-Ecosystem/go-common/pkg/uuidx"
	"github.com/QuizWars-Ecosystem/users-service/internal/metrics"
	"github.com/google/uuid"

	"github.com/QuizWars-Ecosystem/go-common/pkg/abstractions"
	"github.com/QuizWars-Ecosystem/go-common/pkg/jwt"
	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/admin"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/profile"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) SearchUsers(ctx context.Context, request *userspb.SearchUsersRequest) (*userspb.SearchUsersResponse, error) {
	defer metrics.AdminActionsTotalCounter.WithLabelValues("SearchUsers").Inc()

	err := h.jwt.ValidateRoleWithContext(ctx, string(jwt.Admin))
	if err != nil {
		metrics.AdminForbittenActionsTotalCounter.WithLabelValues("SearchUsers", err.Error()).Inc()
		return nil, err
	}

	req, err := abstractions.MakeRequest[admin.SearchFilter](request)
	if err != nil {
		return nil, err
	}

	res, err := h.service.AdminSearchUsers(ctx, req)
	if err != nil {
		return nil, err
	}

	result, err := abstractions.MakeResponse(res)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (h *Handler) GetUserByIdentifier(ctx context.Context, request *userspb.GetUserByIdentifierRequest) (*userspb.UserAdmin, error) {
	defer metrics.AdminActionsTotalCounter.WithLabelValues("GetUserByIdentifier").Inc()

	err := h.jwt.ValidateRoleWithContext(ctx, string(jwt.Admin))
	if err != nil {
		metrics.AdminForbittenActionsTotalCounter.WithLabelValues("GetUserByIdentifier", err.Error()).Inc()
		return nil, err
	}

	var user *profile.UserAdmin

	switch request.Identifier.(type) {
	case *userspb.GetUserByIdentifierRequest_UserId:
		var userID uuid.UUID
		userID, err = uuidx.Parse(request.GetUserId())
		if err != nil {
			return nil, err
		}
		user, err = h.service.AdminGetUserByID(ctx, userID)
	case *userspb.GetUserByIdentifierRequest_Username:
		user, err = h.service.AdminGetUserByUsername(ctx, request.GetUsername())
	case *userspb.GetUserByIdentifierRequest_Email:
		user, err = h.service.AdminGetUserByEmail(ctx, request.GetEmail())
	}

	if err != nil {
		return nil, err
	}

	result, err := abstractions.MakeResponse(user)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (h *Handler) UpdateUserRole(ctx context.Context, request *userspb.UpdateUserRoleRequest) (*emptypb.Empty, error) {
	defer metrics.AdminActionsTotalCounter.WithLabelValues("UpdateUserRole").Inc()

	claims, err := h.jwt.ValidateTokenWithContext(ctx)
	if err != nil {
		metrics.AdminForbittenActionsTotalCounter.WithLabelValues("UpdateUserRole", err.Error()).Inc()
		return nil, err
	}

	if claims.Role != string(jwt.Super) {
		return nil, apperrors.Forbidden(jwt.AuthPermissionDeniedError)
	}

	userID, err := uuidx.Parse(request.GetUserId())
	if err != nil {
		return nil, err
	}

	if err = h.service.AdminUpdateUserRole(ctx, userID, admin.RoleFromGRPCEnum(request.GetRole()).String()); err != nil {
		return nil, err
	}

	return Empty, nil
}

func (h *Handler) BanUser(ctx context.Context, request *userspb.BanUserRequest) (*emptypb.Empty, error) {
	defer metrics.AdminActionsTotalCounter.WithLabelValues("BanUser").Inc()

	err := h.jwt.ValidateRoleWithContext(ctx, string(jwt.Admin))
	if err != nil {
		metrics.AdminForbittenActionsTotalCounter.WithLabelValues("BanUser", err.Error()).Inc()
		return nil, err
	}

	userID, err := uuidx.Parse(request.GetUserId())
	if err != nil {
		return nil, err
	}

	err = h.service.AdminBanUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return Empty, nil
}

func (h *Handler) UnbanUser(ctx context.Context, request *userspb.UnbanUserRequest) (*emptypb.Empty, error) {
	defer metrics.AdminActionsTotalCounter.WithLabelValues("UnbanUser").Inc()

	err := h.jwt.ValidateRoleWithContext(ctx, string(jwt.Admin))
	if err != nil {
		metrics.AdminForbittenActionsTotalCounter.WithLabelValues("UnbanUser", err.Error()).Inc()
		return nil, err
	}

	userID, err := uuidx.Parse(request.GetUserId())
	if err != nil {
		return nil, err
	}

	err = h.service.AdminUnbanUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return Empty, nil
}
