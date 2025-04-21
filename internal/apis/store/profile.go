package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/QuizWars-Ecosystem/users-service/internal/models/profile"
)

func (s *Store) GetProfile(ctx context.Context, userID uuid.UUID) (*profile.Profile, error) {
	return s.db.GetProfile(ctx, userID)
}

func (s *Store) GetUserByID(ctx context.Context, userID uuid.UUID) (*profile.User, error) {
	return s.db.GetUserByID(ctx, userID)
}

func (s *Store) GetUserByUsername(ctx context.Context, username string) (*profile.User, error) {
	return s.db.GetUserByUsername(ctx, username)
}

func (s *Store) UpdateProfile(ctx context.Context, userID uuid.UUID, request *profile.UpdateProfile) error {
	return s.db.UpdateProfile(ctx, userID, request)
}

func (s *Store) UpdateProfileAvatar(ctx context.Context, userID uuid.UUID, avatarID int32) error {
	return s.db.UpdateProfileAvatar(ctx, userID, avatarID)
}

func (s *Store) UpdateProfilePassword(ctx context.Context, userID uuid.UUID, password string) error {
	return s.db.UpdateProfilePassword(ctx, userID, password)
}

func (s *Store) SetProfileRating(ctx context.Context, userID uuid.UUID, rating int32) error {
	return s.db.SetProfileRating(ctx, userID, rating)
}

func (s *Store) SetProfileCoins(ctx context.Context, userID uuid.UUID, coins int64) error {
	return s.db.SetProfileCoins(ctx, userID, coins)
}

func (s *Store) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	return s.db.DeleteProfile(ctx, userID)
}
