package store

import (
	"context"
	"github.com/google/uuid"

	"github.com/QuizWars-Ecosystem/users-service/internal/models/admin"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/profile"
)

func (s *Store) AdminSearchUsers(ctx context.Context, filter *admin.SearchFilter) ([]*profile.UserAdmin, int, error) {
	return s.db.AdminSearchUsers(ctx, filter)
}

func (s *Store) AdminGetUserByID(ctx context.Context, userID uuid.UUID) (*profile.UserAdmin, error) {
	return s.db.AdminGetUserByID(ctx, userID)
}

func (s *Store) AdminGetUserByUsername(ctx context.Context, username string) (*profile.UserAdmin, error) {
	return s.db.AdminGetUserByUsername(ctx, username)
}

func (s *Store) AdminGetUserByEmail(ctx context.Context, email string) (*profile.UserAdmin, error) {
	return s.db.AdminGetUserByEmail(ctx, email)
}

func (s *Store) AdminUpdateUserRole(ctx context.Context, userID uuid.UUID, role string) error {
	return s.db.AdminUpdateUserRole(ctx, userID, role)
}

func (s *Store) AdminBanUser(ctx context.Context, userID uuid.UUID) error {
	return s.db.AdminBanUser(ctx, userID)
}

func (s *Store) AdminUnbanUser(ctx context.Context, userID uuid.UUID) error {
	return s.db.AdminUnbanUser(ctx, userID)
}
