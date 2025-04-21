package store

import (
	"context"

	"github.com/QuizWars-Ecosystem/users-service/internal/models/auth"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/profile"
)

func (s *Store) SaveProfile(ctx context.Context, p *auth.ProfileWithCredentials) (*profile.Profile, error) {
	return s.db.SaveProfile(ctx, p)
}

func (s *Store) GetProfileByUsername(ctx context.Context, username string) (*auth.ProfileWithCredentials, error) {
	return s.db.GetProfileByUsername(ctx, username)
}

func (s *Store) GetProfileByEmail(ctx context.Context, email string) (*auth.ProfileWithCredentials, error) {
	return s.db.GetProfileByEmail(ctx, email)
}

func (s *Store) SetLastLogin(ctx context.Context, userID string) error {
	return s.db.SetLastLogin(ctx, userID)
}
