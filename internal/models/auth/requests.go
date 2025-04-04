package auth

import (
	"github.com/QuizWars-Ecosystem/go-common/pkg/abstractions"
	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/profile"
)

var _ abstractions.Requestable[ProfileWithCredentials, *userspb.RegisterRequest] = (*ProfileWithCredentials)(nil)

func (p ProfileWithCredentials) Request(req *userspb.RegisterRequest) (*ProfileWithCredentials, error) {
	p.Profile = &profile.Profile{
		User: &profile.User{
			AvatarID: req.GetAvatarId(),
			Username: req.GetUsername(),
		},
		Email: req.GetEmail(),
	}

	p.Password = req.GetPassword()

	return &p, nil
}
