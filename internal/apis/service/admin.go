package service

import (
	"context"
	"github.com/google/uuid"

	"github.com/QuizWars-Ecosystem/users-service/internal/models/admin"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/profile"
	"go.uber.org/zap"
)

func (s *Service) AdminSearchUsers(ctx context.Context, filter *admin.SearchFilter) (*admin.SearchUsersResponse, error) {
	users, amount, err := s.store.AdminSearchUsers(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &admin.SearchUsersResponse{
		Users:  users,
		Page:   filter.Offset,
		Size:   filter.Limit,
		Order:  filter.Order,
		Sort:   filter.Sort,
		Amount: int64(amount),
	}, nil
}

func (s *Service) AdminGetUserByID(ctx context.Context, userID uuid.UUID) (*profile.UserAdmin, error) {
	user, err := s.store.AdminGetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) AdminGetUserByUsername(ctx context.Context, username string) (*profile.UserAdmin, error) {
	user, err := s.store.AdminGetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) AdminGetUserByEmail(ctx context.Context, email string) (*profile.UserAdmin, error) {
	user, err := s.store.AdminGetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) AdminUpdateUserRole(ctx context.Context, userID uuid.UUID, role string) error {
	err := s.store.AdminUpdateUserRole(ctx, userID, role)
	if err != nil {
		return err
	}

	s.logger.Info("user role updated", zap.String("user_id", userID.String()))

	return nil
}

func (s *Service) AdminBanUserByID(ctx context.Context, userID uuid.UUID) error {
	err := s.store.AdminBanUser(ctx, userID)
	if err != nil {
		return err
	}

	s.logger.Info("admin banned user", zap.String("user_id", userID.String()))

	return nil
}

func (s *Service) AdminUnbanUserByID(ctx context.Context, userID uuid.UUID) error {
	err := s.store.AdminUnbanUser(ctx, userID)
	if err != nil {
		return err
	}

	s.logger.Info("admin unbanned user", zap.String("user_id", userID.String()))

	return nil
}
