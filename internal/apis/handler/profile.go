package handler

import (
	"context"
	"github.com/QuizWars-Ecosystem/go-common/pkg/abstractions"
	"github.com/QuizWars-Ecosystem/go-common/pkg/uuidx"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/profile"
	"github.com/google/uuid"
	"go.uber.org/zap"

	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) GetProfile(ctx context.Context, request *userspb.GetProfileRequest) (*userspb.GetProfileResponse, error) {
	claims, err := h.jwt.ValidateTokenWithContext(ctx)
	if err != nil {
		return nil, err
	}

	var res *profile.User
	var result *userspb.User

	switch request.Identifier.(type) {
	case *userspb.GetProfileRequest_UserId:
		var userID uuid.UUID
		userID, err = uuidx.Parse(request.GetUserId())
		if err != nil {
			return nil, err
		}

		if claims.UserID == request.GetUserId() {
			var prof *profile.Profile
			prof, err = h.service.GetSelfProfile(ctx, userID)
			if err != nil {
				return nil, err
			}

			var r *userspb.Profile
			r, err = prof.Response()
			if err != nil {
				return nil, err
			}

			return &userspb.GetProfileResponse{
				Data: &userspb.GetProfileResponse_Profile{
					Profile: r,
				},
			}, nil
		}

		res, err = h.service.GetProfileByID(ctx, userID)
		if err != nil {
			return nil, err
		}

		result, err = abstractions.MakeResponse(res)
		if err != nil {
			return nil, err
		}
	case *userspb.GetProfileRequest_Username:
		res, err = h.service.GetProfileByUsername(ctx, request.GetUsername())
		if err != nil {
			return nil, err
		}

		result, err = abstractions.MakeResponse(res)
		if err != nil {
			return nil, err
		}
	}

	return &userspb.GetProfileResponse{
		Data: &userspb.GetProfileResponse_User{
			User: result,
		},
	}, nil
}

func (h *Handler) UpdateProfile(ctx context.Context, request *userspb.UpdateProfileRequest) (*emptypb.Empty, error) {
	err := h.jwt.ValidateUserIDWithContext(ctx, request.GetUserId())
	if err != nil {
		return nil, err
	}

	req, err := abstractions.MakeRequest[profile.UpdateProfile](request)
	if err != nil {
		return nil, err
	}

	userID, err := uuidx.Parse(request.GetUserId())
	if err != nil {
		return nil, err
	}

	err = h.service.UpdateProfile(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	return Empty, nil
}

func (h *Handler) UpdateAvatar(ctx context.Context, request *userspb.UpdateAvatarRequest) (*emptypb.Empty, error) {
	err := h.jwt.ValidateUserIDWithContext(ctx, request.GetUserId())
	if err != nil {
		return nil, err
	}

	userID, err := uuidx.Parse(request.GetUserId())
	if err != nil {
		return nil, err
	}

	err = h.service.UpdateProfileAvatar(ctx, userID, request.GetAvatarId())
	if err != nil {
		return nil, err
	}

	return Empty, nil
}

func (h *Handler) ChangePassword(ctx context.Context, request *userspb.ChangePasswordRequest) (*emptypb.Empty, error) {
	err := h.jwt.ValidateUserIDWithContext(ctx, request.GetUserId())
	if err != nil {
		return nil, err
	}

	userID, err := uuidx.Parse(request.GetUserId())
	if err != nil {
		return nil, err
	}

	err = h.service.UpdateProfilePassword(ctx, userID, request.GetPassword())
	if err != nil {
		return nil, err
	}

	return Empty, nil
}

func (h *Handler) DeleteAccount(ctx context.Context, request *userspb.DeleteAccountRequest) (*emptypb.Empty, error) {
	err := h.jwt.ValidateUserIDWithContext(ctx, request.GetUserId())
	if err != nil {
		return nil, err
	}

	userID, err := uuidx.Parse(request.GetUserId())
	if err != nil {
		return nil, err
	}

	err = h.service.DeleteProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	h.logger.Debug("user deleted", zap.String("id", request.GetUserId()))

	return Empty, nil
}
