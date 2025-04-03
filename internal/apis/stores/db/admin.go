package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Admin struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewAdmin(db *pgxpool.Pool, logger *zap.Logger) *Admin {
	return &Admin{
		db:     db,
		logger: logger,
	}
}
