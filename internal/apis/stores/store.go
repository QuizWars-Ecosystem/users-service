package handlers

import (
	"github.com/QuizWars-Ecosystem/users-service/internal/apis/stores/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Store struct {
	*db.Admin
	*db.Auth
	*db.Profile
	*db.Social
	logger *zap.Logger
}

func NewStore(pool *pgxpool.Pool, logger *zap.Logger) *Store {
	return &Store{
		Admin:   db.NewAdmin(pool, logger),
		Auth:    db.NewAuth(pool, logger),
		Profile: db.NewProfile(pool, logger),
		Social:  db.NewSocial(pool, logger),
		logger:  logger,
	}
}
