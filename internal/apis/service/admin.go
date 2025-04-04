package service

import (
	"context"

	"github.com/QuizWars-Ecosystem/users-service/internal/models/admin"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/profile"
	"go.uber.org/zap"
)

func (s *Service) AdminSearchUsers(ctx context.Context, filter *admin.SearchFilter) (*admin.SearchUsersResponse, error) {
	users, amount, err := s.store.Admin.SearchUsers(ctx, filter)
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

func (s *Service) AdminGetUserByID(ctx context.Context, userID string) (*profile.UserAdmin, error) {
	user, err := s.store.Admin.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) AdminGetUserByUsername(ctx context.Context, username string) (*profile.UserAdmin, error) {
	user, err := s.store.Admin.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) AdminGetUserByEmail(ctx context.Context, email string) (*profile.UserAdmin, error) {
	user, err := s.store.Admin.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) AdminBanUserByID(ctx context.Context, userID string) error {
	err := s.store.Admin.BanUser(ctx, userID)
	if err != nil {
		return err
	}

	s.logger.Info("admin banned user", zap.String("user_id", userID))

	return nil
}

func (s *Service) AdminUnbanUserByID(ctx context.Context, userID string) error {
	err := s.store.Admin.UnbanUser(ctx, userID)
	if err != nil {
		return err
	}

	s.logger.Info("admin unbanned user", zap.String("user_id", userID))

	return nil
}
