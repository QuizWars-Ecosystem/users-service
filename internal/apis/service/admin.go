package service

import (
	"context"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/admin"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/profile"
)

func (s *Service) AdminSearchUsers(ctx context.Context, filter *admin.SearchFilter) (*admin.SearchUsersResponse, error) {

}

func (s *Service) AdminGetUserByID(ctx context.Context, userID string) (*profile.UserAdmin, error) {

}

func (s *Service) AdminGetUserByUsername(ctx context.Context, username string) (*profile.UserAdmin, error) {

}

func (s *Service) AdminGetUserByEmail(ctx context.Context, email string) (*profile.UserAdmin, error) {

}

func (s *Service) AdminBanUserByID(ctx context.Context, userID string) error {

}

func (s *Service) AdminUnbanUserByID(ctx context.Context, userID string) error {

}
