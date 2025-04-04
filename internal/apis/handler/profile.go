package handler

import (
	"context"
	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) GetProfile(ctx context.Context, request *userspb.GetProfileRequest) (*userspb.GetProfileResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) UpdateProfile(ctx context.Context, request *userspb.UpdateProfileRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) UpdateAvatar(ctx context.Context, request *userspb.UpdateAvatarRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) ChangePassword(ctx context.Context, request *userspb.ChangePasswordRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) DeleteAccount(ctx context.Context, request *userspb.DeleteAccountRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
