package profile

import (
	"time"

	"github.com/google/uuid"

	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
)

type Profile struct {
	User  *User
	Email string `json:"email"`
	Coins int64  `json:"coins"`
}

type User struct {
	ID          uuid.UUID  `json:"id"`
	AvatarID    int32      `json:"avatar_id"`
	Username    string     `json:"username"`
	Rating      int32      `json:"rating"`
	CreatedAt   time.Time  `json:"created_at"`
	LastLoginAt *time.Time `json:"last_login_at"`
}

type UserAdmin struct {
	Profile   *Profile
	DeletedAt *time.Time `json:"deleted_at"`
}

type Friend struct {
	User   *User  `json:"user"`
	Status Status `json:"status"`
}

type UpdateProfile struct {
	Username *string `json:"username"`
}

type Status string

func (s Status) String() string {
	return string(s)
}

const (
	Unknown  Status = "unknown"
	Pending  Status = "pending"
	Accepted Status = "accepted"
	Blocked  Status = "blocked"
)

func (s Status) ToGRPCEnum() userspb.Status {
	switch s {
	case Pending:
		return userspb.Status_STATUS_PENDING
	case Accepted:
		return userspb.Status_STATUS_ACCEPTED
	case Blocked:
		return userspb.Status_STATUS_BLOCKED
	default:
		return userspb.Status_STATUS_UNSPECIFIED
	}
}

func statusFromGRPCEnum(status userspb.Status) Status {
	switch status {
	case userspb.Status_STATUS_PENDING:
		return Pending
	case userspb.Status_STATUS_ACCEPTED:
		return Accepted
	case userspb.Status_STATUS_BLOCKED:
		return Blocked
	default:
		return Unknown
	}
}
