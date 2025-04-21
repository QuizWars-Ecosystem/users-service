package store

import (
	"github.com/QuizWars-Ecosystem/users-service/internal/apis/store/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var _ IStore = (*Store)(nil)

type Store struct {
	db     *db.Database
	logger *zap.Logger
}

func NewStore(pool *pgxpool.Pool, logger *zap.Logger) *Store {
	return &Store{
		db:     db.NewDatabase(pool, logger),
		logger: logger,
	}
}
