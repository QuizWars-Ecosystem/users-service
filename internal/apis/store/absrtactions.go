package store

import (
	"context"

	"github.com/google/uuid"

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
	SetLastLogin(ctx context.Context, userID uuid.UUID) error
}

type IProfileStore interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*profile.Profile, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*profile.User, error)
	GetUserByUsername(ctx context.Context, username string) (*profile.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, request *profile.UpdateProfile) error
	UpdateProfileAvatar(ctx context.Context, userID uuid.UUID, avatarID int32) error
	UpdateProfilePassword(ctx context.Context, userID uuid.UUID, password string) error
	SetProfileRating(ctx context.Context, userID uuid.UUID, rating int32) error
	SetProfileCoins(ctx context.Context, userID uuid.UUID, coins int64) error
	DeleteProfile(ctx context.Context, userID uuid.UUID) error
}

type ISocialStore interface {
	AddFriend(ctx context.Context, requesterID, recipientID uuid.UUID) error
	AcceptFriend(ctx context.Context, recipientID, requesterID uuid.UUID) error
	RejectFriend(ctx context.Context, recipientID, requesterID uuid.UUID) error
	RemoveFriend(ctx context.Context, userID, friendID uuid.UUID) error
	GetFriends(ctx context.Context, userID uuid.UUID) ([]*profile.Friend, error)
	BanFriend(ctx context.Context, userID, friendID uuid.UUID) error
	UnbanFriend(ctx context.Context, userID, friendID uuid.UUID) error
}

type IAdminStore interface {
	AdminSearchUsers(ctx context.Context, filter *admin.SearchFilter) ([]*profile.UserAdmin, int, error)
	AdminGetUserByID(ctx context.Context, userID uuid.UUID) (*profile.UserAdmin, error)
	AdminGetUserByUsername(ctx context.Context, username string) (*profile.UserAdmin, error)
	AdminGetUserByEmail(ctx context.Context, email string) (*profile.UserAdmin, error)
	AdminUpdateUserRole(ctx context.Context, userID uuid.UUID, role string) error
	AdminBanUser(ctx context.Context, userID uuid.UUID) error
	AdminUnbanUser(ctx context.Context, userID uuid.UUID) error
}
