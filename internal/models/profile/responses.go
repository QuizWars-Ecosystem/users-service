package profile

import (
	"github.com/QuizWars-Ecosystem/go-common/pkg/abstractions"
	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ abstractions.Responseable[userspb.User] = (*User)(nil)

func (u *User) Response() (*userspb.User, error) {
	var res userspb.User

	res.Id = u.ID
	res.Username = u.Username
	res.AvatarId = u.AvatarID
	res.Rating = u.Rating
	res.CreatedAt = timestamppb.New(u.CreatedAt)

	if u.LastLoginAt != nil {
		res.LastLoginAt = timestamppb.New(*u.LastLoginAt)
	}

	return &res, nil
}

var _ abstractions.Responseable[userspb.Profile] = (*Profile)(nil)

func (p *Profile) Response() (*userspb.Profile, error) {
	var res userspb.Profile

	res.Id = p.User.ID
	res.Username = p.User.Username
	res.AvatarId = p.User.AvatarID
	res.Rating = p.User.Rating
	res.CreatedAt = timestamppb.New(p.User.CreatedAt)

	if p.User.LastLoginAt != nil {
		res.LastLoginAt = timestamppb.New(*p.User.LastLoginAt)
	}

	res.Email = p.Email
	res.Coins = p.Coins

	return &res, nil
}

var _ abstractions.Responseable[userspb.UserAdmin] = (*UserAdmin)(nil)

func (u *UserAdmin) Response() (*userspb.UserAdmin, error) {
	var res userspb.UserAdmin

	res.Id = u.Profile.User.ID
	res.Username = u.Profile.User.Username
	res.AvatarId = u.Profile.User.AvatarID
	res.Rating = u.Profile.User.Rating
	res.CreatedAt = timestamppb.New(u.Profile.User.CreatedAt)

	if u.Profile.User.LastLoginAt != nil {
		res.LastLoginAt = timestamppb.New(*u.Profile.User.LastLoginAt)
	}

	res.Email = u.Profile.Email
	res.Coins = u.Profile.Coins

	if u.DeletedAt != nil {
		res.DeletedAt = timestamppb.New(*u.DeletedAt)
	}

	return &res, nil
}

var _ abstractions.Responseable[userspb.Friend] = (*Friend)(nil)

func (f *Friend) Response() (*userspb.Friend, error) {
	var res userspb.Friend

	user, err := f.User.Response()
	if err != nil {
		return nil, err
	}

	res.User = user
	res.Status = f.Status.ToGRPCEnum()

	return &res, nil
}
