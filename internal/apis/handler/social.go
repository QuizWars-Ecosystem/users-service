package handler

import (
	"context"

	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) AddFriend(ctx context.Context, request *userspb.AddFriendRequest) (*emptypb.Empty, error) {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) AcceptFriend(ctx context.Context, request *userspb.AcceptFriendRequest) (*emptypb.Empty, error) {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) RemoveFriend(ctx context.Context, request *userspb.RemoveFriendRequest) (*emptypb.Empty, error) {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) ListFriends(ctx context.Context, request *userspb.ListFriendsRequest) (*userspb.FriendsList, error) {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) BlockFriend(ctx context.Context, request *userspb.BlockFriendRequest) (*emptypb.Empty, error) {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) UnblockFriend(ctx context.Context, request *userspb.UnblockFriendRequest) (*emptypb.Empty, error) {
	// TODO implement me
	panic("implement me")
}
