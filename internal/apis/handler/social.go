package handler

import (
	"context"

	"github.com/QuizWars-Ecosystem/go-common/pkg/abstractions"

	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) AddFriend(ctx context.Context, request *userspb.AddFriendRequest) (*emptypb.Empty, error) {
	err := h.service.AddFriend(ctx, request.GetRequesterId(), request.GetRecipientId())
	if err != nil {
		return nil, err
	}

	return Empty, nil
}

func (h *Handler) AcceptFriend(ctx context.Context, request *userspb.AcceptFriendRequest) (*emptypb.Empty, error) {
	err := h.service.AcceptFriend(ctx, request.GetRecipientId(), request.GetRequesterId())
	if err != nil {
		return nil, err
	}

	return Empty, nil
}

func (h *Handler) RejectFriend(ctx context.Context, request *userspb.RejectFriendRequest) (*emptypb.Empty, error) {
	err := h.service.RejectFriend(ctx, request.GetRecipientId(), request.GetRequesterId())
	if err != nil {
		return nil, err
	}

	return Empty, nil
}

func (h *Handler) RemoveFriend(ctx context.Context, request *userspb.RemoveFriendRequest) (*emptypb.Empty, error) {
	err := h.jwt.ValidateUserIDToken(request.GetToken(), request.GetRequesterId())
	if err != nil {
		return nil, err
	}

	err = h.service.RemoveFriend(ctx, request.GetRequesterId(), request.GetFriendId())
	if err != nil {
		return nil, err
	}

	return Empty, nil
}

func (h *Handler) ListFriends(ctx context.Context, request *userspb.ListFriendsRequest) (*userspb.FriendsList, error) {
	res, err := h.service.GetFriends(ctx, request.GetUserId())
	if err != nil {
		return nil, err
	}

	friends := make([]*userspb.Friend, len(res))
	for i, f := range res {
		var friend *userspb.Friend

		friend, err = abstractions.MakeResponse(f)
		if err != nil {
			return nil, err
		}

		friends[i] = friend
	}

	return &userspb.FriendsList{
		Friends: friends,
	}, nil
}

func (h *Handler) BlockFriend(ctx context.Context, request *userspb.BlockFriendRequest) (*emptypb.Empty, error) {
	err := h.jwt.ValidateUserIDToken(request.GetToken(), request.GetUserId())
	if err != nil {
		return nil, err
	}

	err = h.service.BlockFriend(ctx, request.GetUserId(), request.GetFriendId())
	if err != nil {
		return nil, err
	}

	return Empty, nil
}

func (h *Handler) UnblockFriend(ctx context.Context, request *userspb.UnblockFriendRequest) (*emptypb.Empty, error) {
	err := h.jwt.ValidateUserIDToken(request.GetToken(), request.GetUserId())
	if err != nil {
		return nil, err
	}

	err = h.service.UnblockFriend(ctx, request.GetUserId(), request.GetFriendId())
	if err != nil {
		return nil, err
	}

	return Empty, nil
}
