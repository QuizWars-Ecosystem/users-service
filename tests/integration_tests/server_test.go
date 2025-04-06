package integration_tests

import (
	"testing"

	test "github.com/QuizWars-Ecosystem/go-common/pkg/testing/server"
	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
	"github.com/QuizWars-Ecosystem/users-service/internal/server"
	"github.com/QuizWars-Ecosystem/users-service/tests/integration_tests/config"
	"github.com/QuizWars-Ecosystem/users-service/tests/integration_tests/modules"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	testCtx := t.Context()
	cfg := config.NewTestConfig()

	prepareInfrastructure(testCtx, t, cfg, runServer)
}

func runServer(t *testing.T, cfg *config.TestConfig) {
	srv, err := server.NewTestServer(t.Context(), cfg.ServiceConfig)
	require.NoError(t, err)

	conn, stop := test.RunServer(t, srv, cfg.ServiceConfig.GRPCPort)
	defer stop()

	authClient := userspb.NewUsersAuthServiceClient(conn)
	socialClient := userspb.NewUsersSocialServiceClient(conn)
	profileClient := userspb.NewUsersProfileServiceClient(conn)
	adminClient := userspb.NewUsersAdminServiceClient(conn)

	modules.AuthServiceTest(t, authClient, cfg)
	modules.SocialServiceTest(t, socialClient, cfg)
	modules.ProfileServiceTest(t, profileClient, cfg)
	modules.AdminServiceTest(t, adminClient, cfg)
}
