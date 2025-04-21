package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Database struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewDatabase(pool *pgxpool.Pool, logger *zap.Logger) *Database {
	return &Database{
		pool:   pool,
		logger: logger,
	}
}
