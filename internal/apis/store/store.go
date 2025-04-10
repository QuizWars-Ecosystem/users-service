package store

import (
	"github.com/QuizWars-Ecosystem/users-service/internal/apis/store/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Store struct {
	Admin   *db.Admin
	Auth    *db.Auth
	Profile *db.Profile
	Social  *db.Social
	logger  *zap.Logger
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
