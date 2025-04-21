package service

import (
	"context"

	"github.com/google/uuid"

	apperrors "github.com/QuizWars-Ecosystem/go-common/pkg/error"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/profile"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) GetSelfProfile(ctx context.Context, userID uuid.UUID) (*profile.Profile, error) {
	prof, err := s.store.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	return prof, nil
}

func (s *Service) GetProfileByID(ctx context.Context, userID uuid.UUID) (*profile.User, error) {
	user, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetProfileByUsername(ctx context.Context, username string) (*profile.User, error) {
	user, err := s.store.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID uuid.UUID, req *profile.UpdateProfile) error {
	err := s.store.UpdateProfile(ctx, userID, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateProfileAvatar(ctx context.Context, userID uuid.UUID, avatarID int32) error {
	err := s.store.UpdateProfileAvatar(ctx, userID, avatarID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateProfilePassword(ctx context.Context, userID uuid.UUID, password string) error {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return apperrors.Internal(err)
	}

	err = s.store.UpdateProfilePassword(ctx, userID, string(passHash))
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	err := s.store.DeleteProfile(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}
