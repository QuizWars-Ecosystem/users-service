package service

import (
	"context"
	"github.com/google/uuid"

	"github.com/QuizWars-Ecosystem/users-service/internal/models/profile"
)

func (s *Service) AddFriend(ctx context.Context, requesterID, recipientID uuid.UUID) error {
	err := s.store.AddFriend(ctx, requesterID, recipientID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) AcceptFriend(ctx context.Context, recipientID, requesterID uuid.UUID) error {
	err := s.store.AcceptFriend(ctx, recipientID, requesterID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) RejectFriend(ctx context.Context, recipientID, requesterID uuid.UUID) error {
	err := s.store.RejectFriend(ctx, recipientID, requesterID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) RemoveFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	err := s.store.RemoveFriend(ctx, userID, friendID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetFriends(ctx context.Context, userID uuid.UUID) ([]*profile.Friend, error) {
	friends, err := s.store.GetFriends(ctx, userID)
	if err != nil {
		return nil, err
	}

	return friends, nil
}

func (s *Service) BlockFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	err := s.store.BanFriend(ctx, userID, friendID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UnblockFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	err := s.store.UnbanFriend(ctx, userID, friendID)
	if err != nil {
		return err
	}

	return nil
}
