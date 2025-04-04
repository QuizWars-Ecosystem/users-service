package handler

import (
	"context"

	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) Register(ctx context.Context, request *userspb.RegisterRequest) (*userspb.RegisterResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) Login(ctx context.Context, request *userspb.LoginRequest) (*userspb.LoginResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) Logout(ctx context.Context, request *userspb.LogoutRequest) (*emptypb.Empty, error) {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) OAuthLogin(ctx context.Context, request *userspb.OAuthLoginRequest) (*userspb.OAuthLoginResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) LinkOAuthProvider(ctx context.Context, request *userspb.LinkOAuthProviderRequest) (*emptypb.Empty, error) {
	// TODO implement me
	panic("implement me")
}
