package service

import (
	"context"
	apperrors "github.com/QuizWars-Ecosystem/go-common/pkg/error"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/profile"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) GetSelfProfile(ctx context.Context, userID string) (*profile.Profile, error) {
	prof, err := s.store.Profile.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	return prof, nil
}

func (s *Service) GetProfileByID(ctx context.Context, userID string) (*profile.User, error) {
	user, err := s.store.Profile.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetProfileByUsername(ctx context.Context, username string) (*profile.User, error) {
	user, err := s.store.Profile.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID string, req *profile.UpdateProfile) error {
	err := s.store.Profile.UpdateProfile(ctx, userID, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateProfileAvatar(ctx context.Context, userID string, avatarID int32) error {
	err := s.store.Profile.UpdateProfileAvatar(ctx, userID, avatarID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateProfilePassword(ctx context.Context, userID string, password string) error {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return apperrors.Internal(err)
	}

	err = s.store.Profile.UpdateProfilePassword(ctx, userID, string(passHash))
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteProfile(ctx context.Context, userID string) error {
	err := s.store.Profile.DeleteProfile(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}
