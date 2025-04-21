package store

import (
	"context"

	"github.com/QuizWars-Ecosystem/users-service/internal/models/admin"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/auth"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/profile"
)

type IStore interface {
	IAuthStore
	IProfileStore
	ISocialStore
	IAdminStore
}

type IAuthStore interface {
	SaveProfile(ctx context.Context, p *auth.ProfileWithCredentials) (*profile.Profile, error)
	GetProfileByUsername(ctx context.Context, username string) (*auth.ProfileWithCredentials, error)
	GetProfileByEmail(ctx context.Context, email string) (*auth.ProfileWithCredentials, error)
	SetLastLogin(ctx context.Context, userID string) error
}

type IProfileStore interface {
	GetProfile(ctx context.Context, userID string) (*profile.Profile, error)
	GetUserByID(ctx context.Context, userID string) (*profile.User, error)
	GetUserByUsername(ctx context.Context, username string) (*profile.User, error)
	UpdateProfile(ctx context.Context, userID string, request *profile.UpdateProfile) error
	UpdateProfileAvatar(ctx context.Context, userID string, avatarID int32) error
	UpdateProfilePassword(ctx context.Context, userID string, password string) error
	SetProfileRating(ctx context.Context, userID string, rating int32) error
	SetProfileCoins(ctx context.Context, userID string, coins int64) error
	DeleteProfile(ctx context.Context, userID string) error
}

type ISocialStore interface {
	AddFriend(ctx context.Context, requesterID string, recipientID string) error
	AcceptFriend(ctx context.Context, recipientID string, requesterID string) error
	RejectFriend(ctx context.Context, recipientID string, requesterID string) error
	RemoveFriend(ctx context.Context, userID string, friendID string) error
	GetFriends(ctx context.Context, userID string) ([]*profile.Friend, error)
	BanFriend(ctx context.Context, userID string, friendID string) error
	UnbanFriend(ctx context.Context, userID string, friendID string) error
}

type IAdminStore interface {
	AdminSearchUsers(ctx context.Context, filter *admin.SearchFilter) ([]*profile.UserAdmin, int, error)
	AdminGetUserByID(ctx context.Context, userID string) (*profile.UserAdmin, error)
	AdminGetUserByUsername(ctx context.Context, username string) (*profile.UserAdmin, error)
	AdminGetUserByEmail(ctx context.Context, email string) (*profile.UserAdmin, error)
	AdminUpdateUserRole(ctx context.Context, userID string, role string) error
	AdminBanUser(ctx context.Context, userID string) error
	AdminUnbanUser(ctx context.Context, userID string) error
}
