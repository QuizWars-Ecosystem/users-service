package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Profile struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewProfile(db *pgxpool.Pool, logger *zap.Logger) *Profile {
	return &Profile{
		db:     db,
		logger: logger,
	}
}
