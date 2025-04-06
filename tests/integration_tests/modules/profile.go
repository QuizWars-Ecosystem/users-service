package modules

import (
	"testing"

	jw "github.com/QuizWars-Ecosystem/go-common/pkg/jwt"
	testerror "github.com/QuizWars-Ecosystem/go-common/pkg/testing/errors"
	"github.com/stretchr/testify/require"

	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
	"github.com/QuizWars-Ecosystem/users-service/tests/integration_tests/config"
)

func ProfileServiceTest(t *testing.T, client userspb.UsersProfileServiceClient, _ *config.TestConfig) {
	ctx := t.Context()

	t.Run("profile.GetProfile: by user id: token not provided", func(t *testing.T) {
		res, err := client.GetProfile(ctx, &userspb.GetProfileRequest{
			Token: "",
			Identifier: &userspb.GetProfileRequest_UserId{
				UserId: lukas.Id,
			},
		})

		require.Error(t, err)
		require.Nil(t, res)
		testerror.RequireForbiddenError(t, err, jw.AuthAccessTokenNotProvidedError)
	})

	t.Run("profile.GetProfile: by user id: invalid token", func(t *testing.T) {
		res, err := client.GetProfile(ctx, &userspb.GetProfileRequest{
			Token: "invalid token",
			Identifier: &userspb.GetProfileRequest_UserId{
				UserId: lukas.Id,
			},
		})

		require.Error(t, err)
		require.Nil(t, res)
		testerror.RequireForbiddenError(t, err, jw.AuthInvalidTokenError)
	})

	t.Run("profile.GetProfile: by user id: not found", func(t *testing.T) {
		testData := "test_id"
		res, err := client.GetProfile(ctx, &userspb.GetProfileRequest{
			Token: johnToken,
			Identifier: &userspb.GetProfileRequest_UserId{
				UserId: testData,
			},
		})

		require.Error(t, err)
		require.Nil(t, res)
		testerror.RequireNotFoundError(t, err, "user", "id", testData)
	})

	t.Run("profile.GetProfile: by user username: not found", func(t *testing.T) {
		testData := "test_username"
		res, err := client.GetProfile(ctx, &userspb.GetProfileRequest{
			Token: johnToken,
			Identifier: &userspb.GetProfileRequest_Username{
				Username: testData,
			},
		})

		require.Error(t, err)
		require.Nil(t, res)
		testerror.RequireNotFoundError(t, err, "user", "username", testData)
	})

	t.Run("profile.GetProfile: self: successful", func(t *testing.T) {
		res, err := client.GetProfile(ctx, &userspb.GetProfileRequest{
			Token: johnToken,
			Identifier: &userspb.GetProfileRequest_UserId{
				UserId: john.Id,
			},
		})

		profile := res.GetProfile()

		require.NoError(t, err)
		require.NotNil(t, profile)
		require.Equal(t, john.Id, profile.Id)
		require.Equal(t, john.Username, profile.Username)
		require.Equal(t, john.Email, profile.Email)
		require.Equal(t, john.AvatarId, profile.AvatarId)
		require.Equal(t, john.Rating, profile.Rating)
		require.Equal(t, john.Coins, profile.Coins)
		require.Equal(t, john.CreatedAt.IsValid(), profile.CreatedAt.IsValid())
	})

	t.Run("profile.GetProfile: by id: successful", func(t *testing.T) {
		res, err := client.GetProfile(ctx, &userspb.GetProfileRequest{
			Token: johnToken,
			Identifier: &userspb.GetProfileRequest_UserId{
				UserId: lukas.Id,
			},
		})

		user := res.GetUser()

		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, lukas.Id, user.Id)
		require.Equal(t, lukas.Username, user.Username)
		require.Equal(t, lukas.AvatarId, user.AvatarId)
		require.Equal(t, lukas.Rating, user.Rating)
		require.Equal(t, lukas.CreatedAt.IsValid(), user.CreatedAt.IsValid())
	})

	t.Run("profile.GetProfile: by username: successful", func(t *testing.T) {
		res, err := client.GetProfile(ctx, &userspb.GetProfileRequest{
			Token: johnToken,
			Identifier: &userspb.GetProfileRequest_Username{
				Username: lukas.Username,
			},
		})

		user := res.GetUser()

		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, lukas.Id, user.Id)
		require.Equal(t, lukas.Username, user.Username)
		require.Equal(t, lukas.AvatarId, user.AvatarId)
		require.Equal(t, lukas.Rating, user.Rating)
		require.Equal(t, lukas.CreatedAt.IsValid(), user.CreatedAt.IsValid())
	})

	t.Run("profile.UpdateProfile: token not provided", func(t *testing.T) {
		testData := "lukas_new"
		_, err := client.UpdateProfile(ctx, &userspb.UpdateProfileRequest{
			Token:    "",
			UserId:   lukas.Id,
			Username: &testData,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthAccessTokenNotProvidedError)
	})

	t.Run("profile.UpdateProfile: invalid token", func(t *testing.T) {
		testData := "lukas_new"
		_, err := client.UpdateProfile(ctx, &userspb.UpdateProfileRequest{
			Token:    "invalid token",
			UserId:   lukas.Id,
			Username: &testData,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthInvalidTokenError)
	})

	t.Run("profile.UpdateProfile: permission denied", func(t *testing.T) {
		testData := "lukas_new"
		_, err := client.UpdateProfile(ctx, &userspb.UpdateProfileRequest{
			Token:    martinToken,
			UserId:   lukas.Id,
			Username: &testData,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthPermissionDeniedError)
	})

	t.Run("profile.UpdateProfile: by admin token: successful", func(t *testing.T) {
		testData := "lukas_new"
		_, err := client.UpdateProfile(ctx, &userspb.UpdateProfileRequest{
			Token:    johnAdminToken,
			UserId:   lukas.Id,
			Username: &testData,
		})

		require.NoError(t, err)
		lukas.Username = testData
	})

	t.Run("profile.UpdateProfile: by self: successful", func(t *testing.T) {
		testData := "lukas_new"
		_, err := client.UpdateProfile(ctx, &userspb.UpdateProfileRequest{
			Token:    lukasToken,
			UserId:   lukas.Id,
			Username: &testData,
		})

		require.NoError(t, err)
		lukas.Username = testData
	})

	t.Run("profile.UpdateAvatar: token not provided", func(t *testing.T) {
		_, err := client.UpdateAvatar(ctx, &userspb.UpdateAvatarRequest{
			Token:    "",
			UserId:   lukas.Id,
			AvatarId: 5,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthAccessTokenNotProvidedError)
	})

	t.Run("profile.UpdateAvatar: invalid token", func(t *testing.T) {
		_, err := client.UpdateAvatar(ctx, &userspb.UpdateAvatarRequest{
			Token:    "invalid token",
			UserId:   lukas.Id,
			AvatarId: 5,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthInvalidTokenError)
	})

	t.Run("profile.UpdateAvatar: permission denied", func(t *testing.T) {
		_, err := client.UpdateAvatar(ctx, &userspb.UpdateAvatarRequest{
			Token:    martinToken,
			UserId:   lukas.Id,
			AvatarId: 5,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthPermissionDeniedError)
	})

	t.Run("profile.UpdateAvatar: not found", func(t *testing.T) {
		testData := "test_id"
		_, err := client.UpdateAvatar(ctx, &userspb.UpdateAvatarRequest{
			Token:    johnAdminToken,
			UserId:   testData,
			AvatarId: 5,
		})

		require.Error(t, err)
		testerror.RequireNotFoundError(t, err, "user", "id", testData)
	})

	t.Run("profile.UpdateAvatar: successful", func(t *testing.T) {
		var testData int32 = 5
		_, err := client.UpdateAvatar(ctx, &userspb.UpdateAvatarRequest{
			Token:    lukasToken,
			UserId:   lukas.Id,
			AvatarId: testData,
		})

		require.NoError(t, err)
		lukas.AvatarId = testData
	})

	t.Run("profile.ChangePassword: token not provided", func(t *testing.T) {
		testData := "new_password"
		_, err := client.ChangePassword(ctx, &userspb.ChangePasswordRequest{
			Token:    "",
			UserId:   lukas.Id,
			Password: testData,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthAccessTokenNotProvidedError)
	})

	t.Run("profile.ChangePassword: invalid token", func(t *testing.T) {
		testData := "new_password"
		_, err := client.ChangePassword(ctx, &userspb.ChangePasswordRequest{
			Token:    "invalid token",
			UserId:   lukas.Id,
			Password: testData,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthInvalidTokenError)
	})

	t.Run("profile.ChangePassword: permission denied", func(t *testing.T) {
		testData := "new_password"
		_, err := client.ChangePassword(ctx, &userspb.ChangePasswordRequest{
			Token:    martinToken,
			UserId:   lukas.Id,
			Password: testData,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthPermissionDeniedError)
	})

	t.Run("profile.ChangePassword: not found", func(t *testing.T) {
		testData := "new_password"
		testDataID := "test_id"
		_, err := client.ChangePassword(ctx, &userspb.ChangePasswordRequest{
			Token:    johnAdminToken,
			UserId:   testDataID,
			Password: testData,
		})

		require.Error(t, err)
		testerror.RequireNotFoundError(t, err, "user", "id", testDataID)
	})

	t.Run("profile.ChangePassword: by admin: successful", func(t *testing.T) {
		testData := "new_password"
		_, err := client.ChangePassword(ctx, &userspb.ChangePasswordRequest{
			Token:    johnAdminToken,
			UserId:   lukas.Id,
			Password: testData,
		})

		require.NoError(t, err)
		lukasPassword = testData
	})

	t.Run("profile.ChangePassword: successful", func(t *testing.T) {
		testData := "new_password"
		_, err := client.ChangePassword(ctx, &userspb.ChangePasswordRequest{
			Token:    martinToken,
			UserId:   martin.Id,
			Password: testData,
		})

		require.NoError(t, err)
		martinPassword = testData
	})

	t.Run("profile.DeleteAccount: token not provided", func(t *testing.T) {
		_, err := client.DeleteAccount(ctx, &userspb.DeleteAccountRequest{
			Token:  "",
			UserId: sonia.Id,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthAccessTokenNotProvidedError)
	})

	t.Run("profile.DeleteAccount: invalid token", func(t *testing.T) {
		_, err := client.DeleteAccount(ctx, &userspb.DeleteAccountRequest{
			Token:  "invalid token",
			UserId: sonia.Id,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthInvalidTokenError)
	})

	t.Run("profile.DeleteAccount: permission denied", func(t *testing.T) {
		_, err := client.DeleteAccount(ctx, &userspb.DeleteAccountRequest{
			Token:  martinToken,
			UserId: sonia.Id,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthPermissionDeniedError)
	})

	t.Run("profile.DeleteAccount: not found", func(t *testing.T) {
		testData := "test_id"
		_, err := client.DeleteAccount(ctx, &userspb.DeleteAccountRequest{
			Token:  johnAdminToken,
			UserId: testData,
		})

		require.Error(t, err)
		testerror.RequireNotFoundError(t, err, "user", "id", testData)
	})

	t.Run("profile.DeleteAccount: by admin: successful", func(t *testing.T) {
		_, err := client.DeleteAccount(ctx, &userspb.DeleteAccountRequest{
			Token:  johnAdminToken,
			UserId: lukas.Id,
		})

		require.NoError(t, err)
	})

	t.Run("profile.DeleteAccount: successful", func(t *testing.T) {
		_, err := client.DeleteAccount(ctx, &userspb.DeleteAccountRequest{
			Token:  martinToken,
			UserId: martin.Id,
		})

		require.NoError(t, err)
	})
}
