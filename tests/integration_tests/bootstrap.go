package integration_tests

import (
	"context"
	"testing"

	"github.com/QuizWars-Ecosystem/go-common/pkg/testing/migrations"

	"github.com/QuizWars-Ecosystem/go-common/pkg/testing/containers"
	"github.com/QuizWars-Ecosystem/users-service/tests/integration_tests/config"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

type runServerFn func(t *testing.T, cfg *config.TestConfig)

func prepareInfrastructure(
	ctx context.Context,
	t *testing.T,
	cfg *config.TestConfig,
	runServerFn runServerFn,
) {
	postgres, err := containers.NewPostgresContainer(ctx, cfg.Postgres)
	require.NoError(t, err)

	defer testcontainers.CleanupContainer(t, postgres)

	postgresUrl, err := postgres.ConnectionString(ctx)
	require.NoError(t, err)

	cfg.ServiceConfig.Postgres.URL = postgresUrl

	migrations.RunMigrations(t, postgresUrl, "../../migrations")

	runServerFn(t, cfg)
}
