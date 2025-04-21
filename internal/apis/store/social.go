package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/QuizWars-Ecosystem/users-service/internal/models/profile"
)

func (s *Store) AddFriend(ctx context.Context, requesterID, recipientID uuid.UUID) error {
	return s.db.AddFriend(ctx, requesterID, recipientID)
}

func (s *Store) AcceptFriend(ctx context.Context, recipientID, requesterID uuid.UUID) error {
	return s.db.AcceptFriend(ctx, recipientID, requesterID)
}

func (s *Store) RejectFriend(ctx context.Context, recipientID, requesterID uuid.UUID) error {
	return s.db.RejectFriend(ctx, recipientID, requesterID)
}

func (s *Store) RemoveFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	return s.db.RemoveFriend(ctx, userID, friendID)
}

func (s *Store) GetFriends(ctx context.Context, userID uuid.UUID) ([]*profile.Friend, error) {
	return s.db.GetFriends(ctx, userID)
}

func (s *Store) BanFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	return s.db.BanFriend(ctx, userID, friendID)
}

func (s *Store) UnbanFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	return s.db.UnbanFriend(ctx, userID, friendID)
}
