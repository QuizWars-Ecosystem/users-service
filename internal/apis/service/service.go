package service

import (
	"github.com/QuizWars-Ecosystem/users-service/internal/apis/store"
	"go.uber.org/zap"
)

type Service struct {
	store  store.IStore
	logger *zap.Logger
}

func NewService(store store.IStore, logger *zap.Logger) *Service {
	return &Service{store: store, logger: logger}
}
