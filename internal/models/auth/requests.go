package auth

import (
	"github.com/QuizWars-Ecosystem/go-common/pkg/abstractions"
	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/profile"
	"github.com/google/uuid"
)

var _ abstractions.Requestable[ProfileWithCredentials, *userspb.RegisterRequest] = (*ProfileWithCredentials)(nil)

func (p ProfileWithCredentials) Request(req *userspb.RegisterRequest) (*ProfileWithCredentials, error) {
	p.Profile = &profile.Profile{
		User: &profile.User{
			ID:       uuid.New().String(),
			AvatarID: req.GetAvatarId(),
			Username: req.GetUsername(),
		},
		Email: req.GetEmail(),
	}

	p.Password = req.GetPassword()

	return &p, nil
}
