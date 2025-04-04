package auth

import "github.com/QuizWars-Ecosystem/users-service/internal/models/profile"

type ProfileWithCredentials struct {
	Profile  *profile.Profile
	Password string
	Role     string
}
