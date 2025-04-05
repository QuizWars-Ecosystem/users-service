package service

import (
	"context"

	"github.com/QuizWars-Ecosystem/users-service/internal/models/profile"
)

func (s *Service) AddFriend(ctx context.Context, userID, friendID string) error {
	err := s.store.Social.AddFriend(ctx, userID, friendID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) AcceptFriend(ctx context.Context, userID, friendID string) error {
	err := s.store.Social.AcceptFriend(ctx, userID, friendID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) RemoveFriend(ctx context.Context, userID, friendID string) error {
	err := s.store.Social.RemoveFriend(ctx, userID, friendID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetFriends(ctx context.Context, userID string) ([]*profile.Friend, error) {
	friends, err := s.store.Social.GetFriends(ctx, userID)
	if err != nil {
		return nil, err
	}

	return friends, nil
}

func (s *Service) BlockFriend(ctx context.Context, userID, friendID string) error {
	err := s.store.Social.BanFriend(ctx, userID, friendID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UnblockFriend(ctx context.Context, userID, friendID string) error {
	err := s.store.Social.UnbanFriend(ctx, userID, friendID)
	if err != nil {
		return err
	}

	return nil
}
