package service

import (
	"context"

	"github.com/google/uuid"

	apperrors "github.com/QuizWars-Ecosystem/go-common/pkg/error"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/auth"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/profile"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Register(ctx context.Context, credits *auth.ProfileWithCredentials) (*profile.Profile, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(credits.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	credits.Password = string(passHash)

	prof, err := s.store.SaveProfile(ctx, credits)
	if err != nil {
		return nil, err
	}

	return prof, nil
}

func (s *Service) LoginByUsername(ctx context.Context, username, password string) (*auth.ProfileWithCredentials, error) {
	prof, err := s.store.GetProfileByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(prof.Password), []byte(password)); err != nil {
		return nil, apperrors.UnauthorizedHidden(err, "wrong credentials")
	}

	return prof, nil
}

func (s *Service) LoginByEmail(ctx context.Context, email, password string) (*auth.ProfileWithCredentials, error) {
	prof, err := s.store.GetProfileByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(prof.Password), []byte(password)); err != nil {
		return nil, apperrors.UnauthorizedHidden(err, "wrong credentials")
	}

	return prof, nil
}

func (s *Service) Logout(ctx context.Context, userID uuid.UUID) error {
	err := s.store.SetLastLogin(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}
