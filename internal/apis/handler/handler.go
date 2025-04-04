package handler

import (
	"github.com/QuizWars-Ecosystem/go-common/pkg/jwt"
	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
	"github.com/QuizWars-Ecosystem/users-service/internal/apis/service"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

var Empty = &emptypb.Empty{}

var _ userspb.UsersAdminServiceServer = (*Handler)(nil)
var _ userspb.UsersAuthServiceServer = (*Handler)(nil)
var _ userspb.UsersProfileServiceServer = (*Handler)(nil)
var _ userspb.UsersSocialServiceServer = (*Handler)(nil)

type Handler struct {
	service *service.Service
	jwt     *jwt.Service
	logger  *zap.Logger
}

func NewHandler(service *service.Service, jwt *jwt.Service, logger *zap.Logger) *Handler {
	return &Handler{service: service, jwt: jwt, logger: logger}
}
