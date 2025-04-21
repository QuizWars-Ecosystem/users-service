package modules

import (
	"context"
	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
)

var (
	martin = &userspb.Profile{
		AvatarId: 2,
		Username: "martin",
		Email:    "martin@mail.com",
	}
	martinPassword = "pass123PASS!"
	martinToken    string
	martinCtx      context.Context
)

var (
	lukas = &userspb.Profile{
		AvatarId: 3,
		Username: "lukas",
		Email:    "lukas@outlook.com",
	}
	lukasPassword = "pass123PASS!"
	lukasToken    string
	lukasCtx      context.Context
)

var (
	sonia = &userspb.Profile{
		AvatarId: 2,
		Username: "sonia",
		Email:    "sonia@outlook.com",
	}
	soniaPassword = "pass123PASS!"
	soniaToken    string
	soniaCtx      context.Context
)

var (
	masha = &userspb.Profile{
		AvatarId: 2,
		Username: "masha",
		Email:    "masha@outlook.com",
	}
	mashaPassword = "pass123PASS!"
	mashaToken    string
	mashaCtx      context.Context
)

var (
	emptyCtx   context.Context
	invalidCtx context.Context
	superCtx   context.Context
)
