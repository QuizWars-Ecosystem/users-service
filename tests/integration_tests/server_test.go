package integration_tests

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
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
	t.Log("Starting server")

	srv, err := server.NewTestServer(t.Context(), cfg.ServiceConfig)
	require.NoError(t, err)

	group := errgroup.Group{}

	group.Go(srv.Start)

	t.Log("Server started")

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", cfg.ServiceConfig.GRPCPort), opts...)
	require.NoError(t, err)

	authService := userspb.NewUsersAuthServiceClient(conn)
	profileClient := userspb.NewUsersProfileServiceClient(conn)
	socialClient := userspb.NewUsersSocialServiceClient(conn)
	adminClient := userspb.NewUsersAdminServiceClient(conn)

	t.Log("Clients created")

	services.AuthServiceTest(t, authService, cfg)
	services.ProfileServiceTest(t, profileClient, cfg)
	services.SocialServiceTest(t, socialClient, cfg)
	services.AdminServiceTest(t, adminClient, cfg)

	err = conn.Close()
	require.NoError(t, err)

	t.Log("Shutting down server")

	stopCtx, cancel := context.WithTimeout(t.Context(), time.Second*10)
	defer cancel()

	err = srv.Shutdown(stopCtx)
	require.NoError(t, err)

	err = group.Wait()
	require.NoError(t, err)

	t.Log("Server shut down")
}
