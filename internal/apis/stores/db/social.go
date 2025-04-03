package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Social struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewSocial(db *pgxpool.Pool, logger *zap.Logger) *Social {
	return &Social{
		db:     db,
		logger: logger,
	}
}
