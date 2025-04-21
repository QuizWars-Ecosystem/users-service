package handler

import (
	"context"
	"github.com/QuizWars-Ecosystem/go-common/pkg/abstractions"
	"github.com/QuizWars-Ecosystem/go-common/pkg/uuidx"

	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) AddFriend(ctx context.Context, request *userspb.AddFriendRequest) (*emptypb.Empty, error) {
	requesterID, err := uuidx.Parse(request.GetRequesterId())
	if err != nil {
		return nil, err
	}

	recipientID, err := uuidx.Parse(request.GetRecipientId())
	if err != nil {
		return nil, err
	}

	err = h.service.AddFriend(ctx, requesterID, recipientID)
	if err != nil {
		return nil, err
	}

	return Empty, nil
}

func (h *Handler) AcceptFriend(ctx context.Context, request *userspb.AcceptFriendRequest) (*emptypb.Empty, error) {
	requesterID, err := uuidx.Parse(request.GetRequesterId())
	if err != nil {
		return nil, err
	}

	recipientID, err := uuidx.Parse(request.GetRecipientId())
	if err != nil {
		return nil, err
	}

	err = h.service.AcceptFriend(ctx, recipientID, requesterID)
	if err != nil {
		return nil, err
	}

	return Empty, nil
}

func (h *Handler) RejectFriend(ctx context.Context, request *userspb.RejectFriendRequest) (*emptypb.Empty, error) {
	requesterID, err := uuidx.Parse(request.GetRequesterId())
	if err != nil {
		return nil, err
	}

	recipientID, err := uuidx.Parse(request.GetRecipientId())
	if err != nil {
		return nil, err
	}

	err = h.service.RejectFriend(ctx, recipientID, requesterID)
	if err != nil {
		return nil, err
	}

	return Empty, nil
}

func (h *Handler) RemoveFriend(ctx context.Context, request *userspb.RemoveFriendRequest) (*emptypb.Empty, error) {
	err := h.jwt.ValidateUserIDWithContext(ctx, request.GetRequesterId())
	if err != nil {
		return nil, err
	}

	requesterID, err := uuidx.Parse(request.GetRequesterId())
	if err != nil {
		return nil, err
	}

	friendID, err := uuidx.Parse(request.GetFriendId())
	if err != nil {
		return nil, err
	}

	err = h.service.RemoveFriend(ctx, requesterID, friendID)
	if err != nil {
		return nil, err
	}

	return Empty, nil
}

func (h *Handler) ListFriends(ctx context.Context, request *userspb.ListFriendsRequest) (*userspb.FriendsList, error) {
	userID, err := uuidx.Parse(request.GetUserId())
	if err != nil {
		return nil, err
	}

	res, err := h.service.GetFriends(ctx, userID)
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
	err := h.jwt.ValidateUserIDWithContext(ctx, request.GetUserId())
	if err != nil {
		return nil, err
	}

	userID, err := uuidx.Parse(request.GetUserId())
	if err != nil {
		return nil, err
	}

	friendID, err := uuidx.Parse(request.GetFriendId())
	if err != nil {
		return nil, err
	}

	err = h.service.BlockFriend(ctx, userID, friendID)
	if err != nil {
		return nil, err
	}

	return Empty, nil
}

func (h *Handler) UnblockFriend(ctx context.Context, request *userspb.UnblockFriendRequest) (*emptypb.Empty, error) {
	err := h.jwt.ValidateUserIDWithContext(ctx, request.GetUserId())
	if err != nil {
		return nil, err
	}

	userID, err := uuidx.Parse(request.GetUserId())
	if err != nil {
		return nil, err
	}

	friendID, err := uuidx.Parse(request.GetFriendId())
	if err != nil {
		return nil, err
	}

	err = h.service.UnblockFriend(ctx, userID, friendID)
	if err != nil {
		return nil, err
	}

	return Empty, nil
}
