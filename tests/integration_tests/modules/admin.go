package modules

import (
	"testing"

	jw "github.com/QuizWars-Ecosystem/go-common/pkg/jwt"
	testerror "github.com/QuizWars-Ecosystem/go-common/pkg/testing/errors"
	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
	"github.com/QuizWars-Ecosystem/users-service/tests/integration_tests/config"
	"github.com/stretchr/testify/require"
)

var (
	page uint64 = 1
	size uint64 = 10
)

func AdminServiceTest(t *testing.T, client userspb.UsersAdminServiceClient, _ *config.TestConfig) {
	ctx := t.Context()

	t.Run("admin.GetUserByIdentifier: access token not provided", func(t *testing.T) {
		_, err := client.GetUserByIdentifier(ctx, &userspb.GetUserByIdentifierRequest{
			Token: "",
			Identifier: &userspb.GetUserByIdentifierRequest_Email{
				Email: martin.Email,
			},
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthAccessTokenNotProvidedError)
	})

	t.Run("admin.GetUserByIdentifier: invalid token", func(t *testing.T) {
		_, err := client.GetUserByIdentifier(ctx, &userspb.GetUserByIdentifierRequest{
			Token: "invalid token",
			Identifier: &userspb.GetUserByIdentifierRequest_Email{
				Email: martin.Email,
			},
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthInvalidTokenError)
	})

	t.Run("admin.GetUserByIdentifier: permission denied", func(t *testing.T) {
		_, err := client.GetUserByIdentifier(ctx, &userspb.GetUserByIdentifierRequest{
			Token: johnToken,
			Identifier: &userspb.GetUserByIdentifierRequest_Email{
				Email: martin.Email,
			},
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthPermissionDeniedError)
	})

	t.Run("admin.GetUserByIdentifier: get by email: not found", func(t *testing.T) {
		testData := "test@mail.com"
		_, err := client.GetUserByIdentifier(ctx, &userspb.GetUserByIdentifierRequest{
			Token: johnAdminToken,
			Identifier: &userspb.GetUserByIdentifierRequest_Email{
				Email: testData,
			},
		})

		require.Error(t, err)
		testerror.RequireNotFoundError(t, err, "user", "email", testData)
	})

	t.Run("admin.GetUserByIdentifier: get by username: not found", func(t *testing.T) {
		testData := "test"
		_, err := client.GetUserByIdentifier(ctx, &userspb.GetUserByIdentifierRequest{
			Token: johnAdminToken,
			Identifier: &userspb.GetUserByIdentifierRequest_Username{
				Username: testData,
			},
		})

		require.Error(t, err)
		testerror.RequireNotFoundError(t, err, "user", "username", testData)
	})

	t.Run("admin.GetUserByIdentifier: get by id: not found", func(t *testing.T) {
		testData := "test_id"
		_, err := client.GetUserByIdentifier(ctx, &userspb.GetUserByIdentifierRequest{
			Token: johnAdminToken,
			Identifier: &userspb.GetUserByIdentifierRequest_UserId{
				UserId: testData,
			},
		})

		require.Error(t, err)
		testerror.RequireNotFoundError(t, err, "user", "id", testData)
	})

	t.Run("admin.GetUserByIdentifier: get by id: successful", func(t *testing.T) {
		res, err := client.GetUserByIdentifier(ctx, &userspb.GetUserByIdentifierRequest{
			Token: johnAdminToken,
			Identifier: &userspb.GetUserByIdentifierRequest_UserId{
				UserId: martin.Id,
			},
		})

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, martin.Id, res.Id)
		require.Equal(t, martin.AvatarId, res.AvatarId)
		require.Equal(t, martin.Username, res.Username)
		require.Equal(t, martin.Email, res.Email)
		require.Equal(t, martin.Rating, res.Rating)
		require.Equal(t, martin.Coins, res.Coins)
		require.Equal(t, martin.CreatedAt, res.CreatedAt)
		require.Equal(t, martin.LastLoginAt, res.LastLoginAt)
	})

	t.Run("admin.GetUserByIdentifier: get by username: successful", func(t *testing.T) {
		res, err := client.GetUserByIdentifier(ctx, &userspb.GetUserByIdentifierRequest{
			Token: johnAdminToken,
			Identifier: &userspb.GetUserByIdentifierRequest_Username{
				Username: martin.Username,
			},
		})

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, martin.Id, res.Id)
		require.Equal(t, martin.AvatarId, res.AvatarId)
		require.Equal(t, martin.Username, res.Username)
		require.Equal(t, martin.Email, res.Email)
		require.Equal(t, martin.Rating, res.Rating)
		require.Equal(t, martin.Coins, res.Coins)
		require.Equal(t, martin.CreatedAt, res.CreatedAt)
		require.Equal(t, martin.LastLoginAt, res.LastLoginAt)
	})

	t.Run("admin.GetUserByIdentifier: get by email: successful", func(t *testing.T) {
		res, err := client.GetUserByIdentifier(ctx, &userspb.GetUserByIdentifierRequest{
			Token: johnAdminToken,
			Identifier: &userspb.GetUserByIdentifierRequest_Email{
				Email: martin.Email,
			},
		})

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, martin.Id, res.Id)
		require.Equal(t, martin.AvatarId, res.AvatarId)
		require.Equal(t, martin.Username, res.Username)
		require.Equal(t, martin.Email, res.Email)
		require.Equal(t, martin.Rating, res.Rating)
		require.Equal(t, martin.Coins, res.Coins)
		require.Equal(t, martin.CreatedAt, res.CreatedAt)
		require.Equal(t, martin.LastLoginAt, res.LastLoginAt)
	})

	t.Run("admin.BanUser: access token not provided", func(t *testing.T) {
		_, err := client.BanUser(ctx, &userspb.BanUserRequest{
			Token:  "",
			UserId: martin.Id,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthAccessTokenNotProvidedError)
	})

	t.Run("admin.BanUser: invalid token", func(t *testing.T) {
		_, err := client.BanUser(ctx, &userspb.BanUserRequest{
			Token:  "invalid token",
			UserId: martin.Id,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthInvalidTokenError)
	})

	t.Run("admin.BanUser: permission denied", func(t *testing.T) {
		_, err := client.BanUser(ctx, &userspb.BanUserRequest{
			Token:  johnToken,
			UserId: martin.Id,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthPermissionDeniedError)
	})

	t.Run("admin.BanUser: not found", func(t *testing.T) {
		testData := "test_id"
		_, err := client.BanUser(ctx, &userspb.BanUserRequest{
			UserId: testData,
			Token:  johnAdminToken,
		})

		require.Error(t, err)
		testerror.RequireNotFoundError(t, err, "user", "id", testData)
	})

	t.Run("admin.BanUser: successful", func(t *testing.T) {
		_, err := client.BanUser(ctx, &userspb.BanUserRequest{
			Token:  johnAdminToken,
			UserId: martin.Id,
		})

		require.NoError(t, err)
	})

	t.Run("admin.GetUserByIdentifier: banned: successful", func(t *testing.T) {
		res, err := client.GetUserByIdentifier(ctx, &userspb.GetUserByIdentifierRequest{
			Token: johnAdminToken,
			Identifier: &userspb.GetUserByIdentifierRequest_UserId{
				UserId: martin.Id,
			},
		})

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, martin.Id, res.Id)
		require.True(t, res.DeletedAt.IsValid())
	})

	t.Run("admin.UnbanUser: access token not provided", func(t *testing.T) {
		_, err := client.UnbanUser(ctx, &userspb.UnbanUserRequest{
			Token:  "",
			UserId: martin.Id,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthAccessTokenNotProvidedError)
	})

	t.Run("admin.UnbanUser: invalid token", func(t *testing.T) {
		_, err := client.UnbanUser(ctx, &userspb.UnbanUserRequest{
			Token:  "invalid token",
			UserId: martin.Id,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthInvalidTokenError)
	})

	t.Run("admin.UnbanUser: permission denied", func(t *testing.T) {
		_, err := client.UnbanUser(ctx, &userspb.UnbanUserRequest{
			Token:  johnToken,
			UserId: martin.Id,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthPermissionDeniedError)
	})

	t.Run("admin.UnbanUser: not found", func(t *testing.T) {
		testData := "test_id"
		_, err := client.UnbanUser(ctx, &userspb.UnbanUserRequest{
			Token:  johnAdminToken,
			UserId: testData,
		})

		require.Error(t, err)
		testerror.RequireNotFoundError(t, err, "user", "id", testData)
	})

	t.Run("admin.UnbanUser: successful", func(t *testing.T) {
		_, err := client.UnbanUser(ctx, &userspb.UnbanUserRequest{
			Token:  johnAdminToken,
			UserId: martin.Id,
		})

		require.NoError(t, err)
	})

	t.Run("admin.GetUserByIdentifier: unbanned: successful", func(t *testing.T) {
		res, err := client.GetUserByIdentifier(ctx, &userspb.GetUserByIdentifierRequest{
			Token: johnAdminToken,
			Identifier: &userspb.GetUserByIdentifierRequest_UserId{
				UserId: martin.Id,
			},
		})

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, martin.Id, res.Id)
		require.False(t, res.DeletedAt.IsValid())
	})

	t.Run("admin.SearchUsers: access token not provided", func(t *testing.T) {
		res, err := client.SearchUsers(ctx, &userspb.SearchUsersRequest{
			Token: "",
		})

		require.Error(t, err)
		require.Nil(t, res)
		testerror.RequireForbiddenError(t, err, jw.AuthAccessTokenNotProvidedError)
	})

	t.Run("admin.SearchUsers: invalid token", func(t *testing.T) {
		res, err := client.SearchUsers(ctx, &userspb.SearchUsersRequest{
			Token: "invalid token",
		})

		require.Error(t, err)
		require.Nil(t, res)
		testerror.RequireForbiddenError(t, err, jw.AuthInvalidTokenError)
	})

	t.Run("admin.SearchUsers: permission denied", func(t *testing.T) {
		res, err := client.SearchUsers(ctx, &userspb.SearchUsersRequest{
			Token: johnToken,
		})

		require.Error(t, err)
		require.Nil(t, res)
		testerror.RequireForbiddenError(t, err, jw.AuthPermissionDeniedError)
	})

	t.Run("admin.SearchUsers: order by username: sort by desc", func(t *testing.T) {
		order := userspb.Order_ORDER_USERNAME
		sort := userspb.Sort_SORT_DESC
		res, err := client.SearchUsers(ctx, &userspb.SearchUsersRequest{
			Token: johnAdminToken,
			Page:  1,
			Size:  10,
			Order: &order,
			Sort:  &sort,
		})

		require.NoError(t, err)
		require.Equal(t, int(res.Amount), len(res.Users))
		for i := len(res.Users) - 2; i >= 0; i-- {
			require.GreaterOrEqual(t, res.Users[i].Username, res.Users[i+1].Username)
		}
	})

	t.Run("admin.SearchUsers: order by username: sort by  asc", func(t *testing.T) {
		order := userspb.Order_ORDER_USERNAME
		sort := userspb.Sort_SORT_ASC
		res, err := client.SearchUsers(ctx, &userspb.SearchUsersRequest{
			Token: johnAdminToken,
			Page:  1,
			Size:  10,
			Order: &order,
			Sort:  &sort,
		})

		require.NoError(t, err)
		require.Equal(t, int(res.Amount), len(res.Users))
		for i := len(res.Users) - 2; i >= 0; i-- {
			require.LessOrEqual(t, res.Users[i].Username, res.Users[i+1].Username)
		}
	})

	t.Run("admin.SearchUsers: filter by rating: not found", func(t *testing.T) {
		order := userspb.Order_ORDER_USERNAME
		sort := userspb.Sort_SORT_DESC
		res, err := client.SearchUsers(ctx, &userspb.SearchUsersRequest{
			Token: johnAdminToken,
			Page:  1,
			Size:  10,
			Order: &order,
			Sort:  &sort,
			UserRating: &userspb.RatingFiler{
				From: 100,
				To:   100,
			},
		})

		require.Nil(t, err)
		require.Equal(t, 0, len(res.Users))
	})

	t.Run("admin.SearchUsers: filter by rating: successful", func(t *testing.T) {
		order := userspb.Order_ORDER_USERNAME
		sort := userspb.Sort_SORT_DESC
		res, err := client.SearchUsers(ctx, &userspb.SearchUsersRequest{
			Token: johnAdminToken,
			Page:  1,
			Size:  10,
			Order: &order,
			Sort:  &sort,
			UserRating: &userspb.RatingFiler{
				From: 0,
				To:   100,
			},
		})

		require.NoError(t, err)
		require.Equal(t, int(res.Amount), len(res.Users))
		for _, u := range res.Users {
			require.True(t, u.Rating >= 0 && u.Rating <= 100)
		}
	})

	t.Run("admin.SearchUsers: filter by coins: not found", func(t *testing.T) {
		order := userspb.Order_ORDER_USERNAME
		sort := userspb.Sort_SORT_DESC
		res, err := client.SearchUsers(ctx, &userspb.SearchUsersRequest{
			Token: johnAdminToken,
			Page:  1,
			Size:  10,
			Order: &order,
			Sort:  &sort,
			UserCoins: &userspb.CoinsFiler{
				From: 100,
				To:   100,
			},
		})

		require.Nil(t, err)
		require.Equal(t, 0, len(res.Users))
	})

	t.Run("admin.SearchUsers: filter by coins: successful", func(t *testing.T) {
		order := userspb.Order_ORDER_USERNAME
		sort := userspb.Sort_SORT_DESC
		res, err := client.SearchUsers(ctx, &userspb.SearchUsersRequest{
			Token: johnAdminToken,
			Page:  page,
			Size:  size,
			Order: &order,
			Sort:  &sort,
			UserCoins: &userspb.CoinsFiler{
				From: 0,
				To:   100,
			},
		})

		require.NoError(t, err)
		require.Equal(t, int(res.Amount), len(res.Users))
		for _, u := range res.Users {
			require.True(t, u.Coins >= 0 && u.Coins <= 100)
		}
	})
}
