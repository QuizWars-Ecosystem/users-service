package profile

import (
	"github.com/QuizWars-Ecosystem/go-common/pkg/abstractions"
	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
)

var _ abstractions.Requestable[User, *userspb.User] = (*User)(nil)

func (u *User) Request(req *userspb.User) (*User, error) {
	u.ID = req.GetId()
	u.AvatarID = req.GetAvatarId()
	u.Username = req.GetUsername()
	u.Rating = req.GetRating()
	u.CreatedAt = req.CreatedAt.AsTime()

	if req.LastLoginAt != nil {
		time := req.LastLoginAt.AsTime()
		u.LastLoginAt = &time
	}

	return u, nil
}

var _ abstractions.Requestable[Profile, *userspb.Profile] = (*Profile)(nil)

func (p *Profile) Request(req *userspb.Profile) (*Profile, error) {
	p.User = &User{
		ID:        req.GetId(),
		Username:  req.GetUsername(),
		AvatarID:  req.GetAvatarId(),
		Rating:    req.GetRating(),
		CreatedAt: req.CreatedAt.AsTime(),
	}

	if req.LastLoginAt != nil {
		time := req.LastLoginAt.AsTime()
		p.User.LastLoginAt = &time
	}

	p.Email = req.GetEmail()
	p.Coins = req.GetCoins()

	return p, nil
}

var _ abstractions.Requestable[UserAdmin, *userspb.UserAdmin] = (*UserAdmin)(nil)

func (u *UserAdmin) Request(req *userspb.UserAdmin) (*UserAdmin, error) {
	u.Profile.User.ID = req.GetId()
	u.Profile.User.AvatarID = req.GetAvatarId()
	u.Profile.User.Username = req.GetUsername()
	u.Profile.User.Rating = req.GetRating()
	u.Profile.User.CreatedAt = req.CreatedAt.AsTime()

	if req.LastLoginAt != nil {
		time := req.LastLoginAt.AsTime()
		u.Profile.User.LastLoginAt = &time
	}

	u.Profile.Email = req.GetEmail()
	u.Profile.Coins = req.GetCoins()

	if req.DeletedAt != nil {
		time := req.DeletedAt.AsTime()
		u.DeletedAt = &time
	}

	return u, nil
}

var _ abstractions.Requestable[Friend, *userspb.Friend] = (*Friend)(nil)

func (f *Friend) Request(req *userspb.Friend) (*Friend, error) {
	f.User = &User{
		ID:        req.GetUser().GetId(),
		Username:  req.GetUser().GetUsername(),
		AvatarID:  req.GetUser().GetAvatarId(),
		Rating:    req.GetUser().GetRating(),
		CreatedAt: req.GetUser().CreatedAt.AsTime(),
	}

	if req.GetUser().GetLastLoginAt() != nil {
		time := req.GetUser().GetLastLoginAt().AsTime()
		f.User.LastLoginAt = &time
	}

	f.Status = statusFromGRPCEnum(req.GetStatus())

	return f, nil
}
