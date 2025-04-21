package modules

import (
	"context"
	"testing"

	"github.com/google/uuid"

	jw "github.com/QuizWars-Ecosystem/go-common/pkg/jwt"

	testerror "github.com/QuizWars-Ecosystem/go-common/pkg/testing/errors"
	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
	"github.com/QuizWars-Ecosystem/users-service/tests/integration_tests/config"
	"github.com/stretchr/testify/require"
)

var jwt *jw.Service

func AuthServiceTest(t *testing.T, client userspb.UsersAuthServiceClient, cfg *config.TestConfig) {
	ctx := t.Context()

	jwt = jw.NewService(cfg.ServiceConfig.JWT)
	emptyCtx = jwt.SetTokenInContext(ctx, "")
	invalidCtx = jwt.SetTokenInContext(ctx, "invalid token")
	superCtx, _ = jwt.GenerateTokenWithContext(ctx, uuid.New().String(), string(jw.Super))

	t.Run("auth.Register: successful", func(t *testing.T) {
		res, err := client.Register(ctx, &userspb.RegisterRequest{
			AvatarId: john.AvatarId,
			Username: john.Username,
			Email:    john.Email,
			Password: johnPassword,
		})

		require.NoError(t, err)

		profile := res.GetProfile()

		require.NotEqual(t, "", res.GetToken())
		require.NotEqual(t, "", profile.GetId())
		require.Equal(t, john.AvatarId, profile.GetAvatarId())
		require.Equal(t, john.Username, profile.GetUsername())
		require.Equal(t, john.Email, profile.GetEmail())

		claims, err := jwt.ValidateToken(res.GetToken())
		require.NoError(t, err)

		require.Equal(t, profile.GetId(), claims.UserID)
		require.Equal(t, string(jw.User), claims.Role)

		johnToken = res.GetToken()
		john = profile
		johnCtx = jwt.SetTokenInContext(ctx, johnToken)
	})

	t.Run("auth.Register: username already taken", func(t *testing.T) {
		_, err := client.Register(ctx, &userspb.RegisterRequest{
			AvatarId: john.AvatarId,
			Username: john.Username,
			Email:    "test@test.com",
			Password: johnPassword,
		})

		require.Error(t, err)
		testerror.RequireAlreadyExistsError(t, err, "user", "username", john.Username)
	})

	t.Run("auth.Register: email already taken", func(t *testing.T) {
		_, err := client.Register(ctx, &userspb.RegisterRequest{
			AvatarId: john.AvatarId,
			Username: "test",
			Email:    john.Email,
			Password: johnPassword,
		})

		require.Error(t, err)
		testerror.RequireAlreadyExistsError(t, err, "user", "email", john.Email)
	})

	t.Run("auth.Login: login by username: nof found", func(t *testing.T) {
		_, err := client.Login(ctx, &userspb.LoginRequest{
			Identifier: &userspb.LoginRequest_Username{
				Username: martin.Username,
			},
			Password: martinPassword,
		})

		require.Error(t, err)
		testerror.RequireNotFoundError(t, err, "user", "username", martin.Username)
	})

	t.Run("auth.Login: login by email: nof found", func(t *testing.T) {
		_, err := client.Login(ctx, &userspb.LoginRequest{
			Identifier: &userspb.LoginRequest_Email{
				Email: martin.Email,
			},
			Password: martinPassword,
		})

		require.Error(t, err)
		testerror.RequireNotFoundError(t, err, "user", "email", martin.Email)
	})

	t.Run("auth.Logout: not found", func(t *testing.T) {
		_, err := client.Logout(ctx, &userspb.LogoutRequest{
			UserId: martin.Id,
		})

		require.Error(t, err)
		testerror.RequireNotFoundError(t, err, "user", "id", martin.Id)
	})

	t.Run("auth.Logout: successful", func(t *testing.T) {
		_, err := client.Logout(ctx, &userspb.LogoutRequest{
			UserId: john.Id,
		})

		require.NoError(t, err)
	})

	t.Run("auth.Register: successful", func(t *testing.T) {
		res, err := client.Register(ctx, &userspb.RegisterRequest{
			AvatarId: martin.AvatarId,
			Username: martin.Username,
			Email:    martin.Email,
			Password: martinPassword,
		})

		require.NoError(t, err)

		profile := res.GetProfile()

		require.NotEqual(t, "", res.GetToken())
		require.NotEqual(t, "", profile.GetId())
		require.Equal(t, martin.AvatarId, profile.GetAvatarId())
		require.Equal(t, martin.Username, profile.GetUsername())
		require.Equal(t, martin.Email, profile.GetEmail())

		claims, err := jwt.ValidateToken(res.GetToken())
		require.NoError(t, err)

		require.Equal(t, profile.GetId(), claims.UserID)
		require.Equal(t, string(jw.User), claims.Role)

		martinToken = res.GetToken()
		martin = profile
		martinCtx = jwt.SetTokenInContext(ctx, martinToken)
	})

	t.Run("auth.Login: login by email: successful", func(t *testing.T) {
		res, err := client.Login(ctx, &userspb.LoginRequest{
			Identifier: &userspb.LoginRequest_Email{
				Email: martin.Email,
			},
			Password: martinPassword,
		})

		require.NoError(t, err)

		profile := res.GetProfile()

		require.NotEqual(t, "", res.GetToken())
		require.NotEqual(t, "", profile.GetId())
		require.Equal(t, martin.AvatarId, profile.GetAvatarId())
		require.Equal(t, martin.Username, profile.GetUsername())
		require.Equal(t, martin.Email, profile.GetEmail())

		claims, err := jwt.ValidateToken(res.GetToken())
		require.NoError(t, err)

		require.Equal(t, profile.GetId(), claims.UserID)
		require.Equal(t, string(jw.User), claims.Role)

		martinToken = res.GetToken()
		martin = profile
		martinCtx = jwt.SetTokenInContext(ctx, martinToken)
	})

	t.Run("auth.Login: login by username: successful", func(t *testing.T) {
		res, err := client.Login(ctx, &userspb.LoginRequest{
			Identifier: &userspb.LoginRequest_Username{
				Username: john.Username,
			},
			Password: johnPassword,
		})

		require.NoError(t, err)

		profile := res.GetProfile()

		require.NotEqual(t, "", res.GetToken())
		require.NotEqual(t, "", profile.GetId())
		require.Equal(t, john.AvatarId, profile.GetAvatarId())
		require.Equal(t, john.Username, profile.GetUsername())
		require.Equal(t, john.Email, profile.GetEmail())

		claims, err := jwt.ValidateToken(res.GetToken())
		require.NoError(t, err)

		require.Equal(t, profile.GetId(), claims.UserID)
		require.Equal(t, string(jw.User), claims.Role)

		johnToken = res.GetToken()
		john = profile
		johnCtx = jwt.SetTokenInContext(ctx, johnToken)
	})

	t.Run("auth.Register: successful", func(t *testing.T) {
		res, err := client.Register(ctx, &userspb.RegisterRequest{
			AvatarId: lukas.AvatarId,
			Username: lukas.Username,
			Email:    lukas.Email,
			Password: lukasPassword,
		})

		require.NoError(t, err)

		profile := res.GetProfile()

		require.NotEqual(t, "", res.GetToken())
		require.NotEqual(t, "", profile.GetId())
		require.Equal(t, lukas.AvatarId, profile.GetAvatarId())
		require.Equal(t, lukas.Username, profile.GetUsername())
		require.Equal(t, lukas.Email, profile.GetEmail())

		claims, err := jwt.ValidateToken(res.GetToken())
		require.NoError(t, err)

		require.Equal(t, profile.GetId(), claims.UserID)
		require.Equal(t, string(jw.User), claims.Role)

		lukasToken = res.GetToken()
		lukas = profile
		lukasCtx = jwt.SetTokenInContext(ctx, lukasToken)
	})

	t.Run("auth.Register: by list: successful", func(t *testing.T) {
		reqs := []struct {
			data    *userspb.RegisterRequest
			ownerFn func(data *userspb.Profile)
			token   *string
			ctxFn   func(ctx context.Context)
		}{
			{
				data: &userspb.RegisterRequest{
					AvatarId: sonia.AvatarId,
					Username: sonia.Username,
					Email:    sonia.Email,
					Password: soniaPassword,
				},
				ownerFn: func(data *userspb.Profile) {
					sonia = data
				},
				token: &soniaToken,
				ctxFn: func(ctx context.Context) {
					soniaCtx = ctx
				},
			},
			{
				data: &userspb.RegisterRequest{
					AvatarId: masha.AvatarId,
					Username: masha.Username,
					Email:    masha.Email,
					Password: mashaPassword,
				},
				ownerFn: func(data *userspb.Profile) {
					masha = data
				},
				token: &mashaToken,
				ctxFn: func(ctx context.Context) {
					mashaCtx = ctx
				},
			},
		}

		for _, req := range reqs {
			res, err := client.Register(ctx, req.data)

			require.NoError(t, err)

			profile := res.GetProfile()

			require.NotEqual(t, "", res.GetToken())
			require.NotEqual(t, "", profile.GetId())
			require.Equal(t, req.data.GetUsername(), profile.GetUsername())
			require.Equal(t, req.data.GetEmail(), profile.GetEmail())
			require.Equal(t, req.data.GetAvatarId(), profile.GetAvatarId())

			token := res.GetToken()

			req.token = &token
			req.ownerFn(profile)
			req.ctxFn(jwt.SetTokenInContext(ctx, *req.token))
		}
	})

	var err error
	johnAdminToken, err = jwt.GenerateToken(john.GetId(), "admin")
	require.NoError(t, err)
	johnAdminCtx = jwt.SetTokenInContext(ctx, johnAdminToken)
}

var (
	john = &userspb.Profile{
		AvatarId: 1,
		Username: "john",
		Email:    "john@gmail.com",
	}
	johnPassword   = "pass123PASS!"
	johnToken      string
	johnCtx        context.Context
	johnAdminToken string
	johnAdminCtx   context.Context
)
