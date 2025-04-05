package modules

import (
	"testing"

	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
	"github.com/QuizWars-Ecosystem/users-service/tests/integration_tests/config"
	"github.com/stretchr/testify/require"
)

func AuthServiceTest(t *testing.T, client userspb.UsersAuthServiceClient, _ *config.TestConfig) {
	ctx := t.Context()

	t.Run("auth.Register: email is required", func(t *testing.T) {
		_, err := client.Register(ctx, &userspb.RegisterRequest{
			AvatarId: 1,
			Username: "test_1",
			Email:    "test_1@test.com",
			Password: "pass123PASS!",
		})

		require.NoError(t, err)
	})
}
