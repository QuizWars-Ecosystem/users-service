package handler

import (
	"context"

	"github.com/QuizWars-Ecosystem/go-common/pkg/abstractions"
	"github.com/QuizWars-Ecosystem/go-common/pkg/jwt"
	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/admin"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/profile"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) SearchUsers(ctx context.Context, request *userspb.SearchUsersRequest) (*userspb.SearchUsersResponse, error) {
	err := h.jwt.ValidateRoleWithContext(ctx, string(jwt.Admin))
	if err != nil {
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
	err := h.jwt.ValidateRoleWithContext(ctx, string(jwt.Admin))
	if err != nil {
		return nil, err
	}

	var user *profile.UserAdmin

	switch request.Identifier.(type) {
	case *userspb.GetUserByIdentifierRequest_UserId:
		user, err = h.service.AdminGetUserByID(ctx, request.GetUserId())
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

func (h *Handler) BanUser(ctx context.Context, request *userspb.BanUserRequest) (*emptypb.Empty, error) {
	err := h.jwt.ValidateRoleWithContext(ctx, string(jwt.Admin))
	if err != nil {
		return nil, err
	}

	err = h.service.AdminBanUserByID(ctx, request.GetUserId())
	if err != nil {
		return nil, err
	}

	return Empty, nil
}

func (h *Handler) UnbanUser(ctx context.Context, request *userspb.UnbanUserRequest) (*emptypb.Empty, error) {
	err := h.jwt.ValidateRoleWithContext(ctx, string(jwt.Admin))
	if err != nil {
		return nil, err
	}

	err = h.service.AdminUnbanUserByID(ctx, request.GetUserId())
	if err != nil {
		return nil, err
	}

	return Empty, nil
}
