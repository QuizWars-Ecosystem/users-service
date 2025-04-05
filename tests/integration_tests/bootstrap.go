package integration_tests

import (
	"context"
	"github.com/QuizWars-Ecosystem/users-service/tests/integration_tests/config"
	"github.com/QuizWars-Ecosystem/users-service/tests/integration_tests/containers"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/tern/migrate"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"testing"
)

type runServerFn func(t *testing.T, cfg *config.TestConfig)

func prepareInfrastructure(
	ctx context.Context,
	t *testing.T,
	cfg *config.TestConfig,
	runServerFn runServerFn,
) {
	t.Log("Prepare infrastructure...")

	t.Log("Initialize postgres container...")
	postgres, err := containers.NewPostgres(ctx, cfg.Postgres)
	require.NoError(t, err)

	defer testcontainers.CleanupContainer(t, postgres)

	postgresUrl, err := postgres.ConnectionString(ctx)
	require.NoError(t, err)

	cfg.ServiceConfig.Postgres.URL = postgresUrl

	t.Log("Running postgres migrations...")
	runMigrations(t, postgresUrl)

	runServerFn(t, cfg)
}

func runMigrations(t *testing.T, pgConnString string) {
	conn, err := pgx.Connect(t.Context(), pgConnString)
	require.NoError(t, err)

	migrator, err := migrate.NewMigrator(t.Context(), conn, "migrations")
	require.NoError(t, err)

	err = migrator.LoadMigrations("../../migrations")
	require.NoError(t, err)

	err = migrator.Migrate(t.Context())
	require.NoError(t, err)
}
