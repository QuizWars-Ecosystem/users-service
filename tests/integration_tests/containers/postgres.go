package containers

import (
	"context"
	"fmt"
	"time"

	"github.com/QuizWars-Ecosystem/users-service/tests/integration_tests/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewPostgres(ctx context.Context, cfg *config.PostgresConfig) (*postgres.PostgresContainer, error) {
	return postgres.Run(
		ctx,
		cfg.Image,
		postgres.WithDatabase(cfg.DBName),
		postgres.WithUsername(cfg.Username),
		postgres.WithPassword(cfg.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(time.Second*5),
		),
		testcontainers.WithHostPortAccess([]int{cfg.Port}...),
	)
}

func BuildPostgresURL(username, password, host, port, db string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, db)
}
