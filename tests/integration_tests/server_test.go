package integration_tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
	"github.com/QuizWars-Ecosystem/users-service/internal/server"
	"github.com/QuizWars-Ecosystem/users-service/tests/integration_tests/config"
	"github.com/QuizWars-Ecosystem/users-service/tests/integration_tests/services"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Test(t *testing.T) {
	testCtx := t.Context()
	cfg := config.NewTestConfig()

	prepareInfrastructure(testCtx, t, cfg, runServer)
}

func runServer(t *testing.T, cfg *config.TestConfig) {
	srv, err := server.NewTestServer(t.Context(), cfg.ServiceConfig)
	require.NoError(t, err)

	go func() {
		err = srv.Start()
		require.NoError(t, err)
	}()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", cfg.ServiceConfig.GRPCPort), opts...)
	require.NoError(t, err)

	authService := userspb.NewUsersAuthServiceClient(conn)
	profileClient := userspb.NewUsersProfileServiceClient(conn)
	socialClient := userspb.NewUsersSocialServiceClient(conn)
	adminClient := userspb.NewUsersAdminServiceClient(conn)

	services.AuthServiceTest(t, authService, cfg)
	services.ProfileServiceTest(t, profileClient, cfg)
	services.SocialServiceTest(t, socialClient, cfg)
	services.AdminServiceTest(t, adminClient, cfg)

	err = conn.Close()
	require.NoError(t, err)

	stopCtx, cancel := context.WithTimeout(t.Context(), time.Second*10)
	defer cancel()

	err = srv.Shutdown(stopCtx)
	require.NoError(t, err)
}
