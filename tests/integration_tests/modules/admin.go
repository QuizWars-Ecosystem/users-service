package modules

import (
	"testing"

	"github.com/google/uuid"

	jw "github.com/QuizWars-Ecosystem/go-common/pkg/jwt"
	testerror "github.com/QuizWars-Ecosystem/go-common/pkg/testing/errors"
	usersv1 "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
	"github.com/QuizWars-Ecosystem/users-service/tests/integration_tests/config"
	"github.com/stretchr/testify/require"
)

var (
	page uint64 = 1
	size uint64 = 10
)

func AdminServiceTest(t *testing.T, client usersv1.UsersAdminServiceClient, _ *config.TestConfig) {
	t.Run("admin.GetUserByIdentifier: access token not provided", func(t *testing.T) {
		_, err := client.GetUserByIdentifier(emptyCtx, &usersv1.GetUserByIdentifierRequest{
			Identifier: &usersv1.GetUserByIdentifierRequest_Email{
				Email: martin.Email,
			},
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthAccessTokenNotProvidedError)
	})

	t.Run("admin.GetUserByIdentifier: invalid token", func(t *testing.T) {
		_, err := client.GetUserByIdentifier(invalidCtx, &usersv1.GetUserByIdentifierRequest{
			Identifier: &usersv1.GetUserByIdentifierRequest_Email{
				Email: martin.Email,
			},
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthInvalidTokenError)
	})

	t.Run("admin.GetUserByIdentifier: permission denied", func(t *testing.T) {
		_, err := client.GetUserByIdentifier(johnCtx, &usersv1.GetUserByIdentifierRequest{
			Identifier: &usersv1.GetUserByIdentifierRequest_Email{
				Email: martin.Email,
			},
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthPermissionDeniedError)
	})

	t.Run("admin.GetUserByIdentifier: get by email: not found", func(t *testing.T) {
		testData := "test@mail.com"
		_, err := client.GetUserByIdentifier(johnAdminCtx, &usersv1.GetUserByIdentifierRequest{
			Identifier: &usersv1.GetUserByIdentifierRequest_Email{
				Email: testData,
			},
		})

		require.Error(t, err)
		testerror.RequireNotFoundError(t, err, "user", "email", testData)
	})

	t.Run("admin.GetUserByIdentifier: get by username: not found", func(t *testing.T) {
		testData := "test"
		_, err := client.GetUserByIdentifier(johnAdminCtx, &usersv1.GetUserByIdentifierRequest{
			Identifier: &usersv1.GetUserByIdentifierRequest_Username{
				Username: testData,
			},
		})

		require.Error(t, err)
		testerror.RequireNotFoundError(t, err, "user", "username", testData)
	})

	t.Run("admin.GetUserByIdentifier: get by id: not found", func(t *testing.T) {
		testData := "test_id"
		_, err := client.GetUserByIdentifier(johnAdminCtx, &usersv1.GetUserByIdentifierRequest{
			Identifier: &usersv1.GetUserByIdentifierRequest_UserId{
				UserId: testData,
			},
		})

		require.Error(t, err)
		testerror.RequireNotFoundError(t, err, "user", "id", testData)
	})

	t.Run("admin.GetUserByIdentifier: get by id: successful", func(t *testing.T) {
		res, err := client.GetUserByIdentifier(johnAdminCtx, &usersv1.GetUserByIdentifierRequest{
			Identifier: &usersv1.GetUserByIdentifierRequest_UserId{
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
		res, err := client.GetUserByIdentifier(johnAdminCtx, &usersv1.GetUserByIdentifierRequest{
			Identifier: &usersv1.GetUserByIdentifierRequest_Username{
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
		res, err := client.GetUserByIdentifier(johnAdminCtx, &usersv1.GetUserByIdentifierRequest{
			Identifier: &usersv1.GetUserByIdentifierRequest_Email{
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

	t.Run("admin.UpdateUserRole: access token not provided", func(t *testing.T) {
		_, err := client.UpdateUserRole(emptyCtx, &usersv1.UpdateUserRoleRequest{})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthAccessTokenNotProvidedError)
	})

	t.Run("admin.UpdateUserRole: invalid token", func(t *testing.T) {
		_, err := client.UpdateUserRole(invalidCtx, &usersv1.UpdateUserRoleRequest{})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthInvalidTokenError)
	})

	t.Run("admin.UpdateUserRole: user: permission denied", func(t *testing.T) {
		_, err := client.UpdateUserRole(mashaCtx, &usersv1.UpdateUserRoleRequest{})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthPermissionDeniedError)
	})

	t.Run("admin.UpdateUserRole: admin: permission denied", func(t *testing.T) {
		_, err := client.UpdateUserRole(johnAdminCtx, &usersv1.UpdateUserRoleRequest{})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthPermissionDeniedError)
	})

	t.Run("admin.UpdateUserRole: not found", func(t *testing.T) {
		testID := uuid.New().String()
		_, err := client.UpdateUserRole(superCtx, &usersv1.UpdateUserRoleRequest{
			UserId: testID,
			Role:   usersv1.Role_ROLE_USER,
		})

		require.Error(t, err)
		testerror.RequireNotFoundError(t, err, "user", "id", testID)
	})

	t.Run("admin.UpdateUserRole: successful", func(t *testing.T) {
		_, err := client.UpdateUserRole(superCtx, &usersv1.UpdateUserRoleRequest{
			UserId: masha.GetId(),
			Role:   usersv1.Role_ROLE_USER,
		})

		require.NoError(t, err)
	})

	t.Run("admin.BanUser: access token not provided", func(t *testing.T) {
		_, err := client.BanUser(emptyCtx, &usersv1.BanUserRequest{
			UserId: martin.Id,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthAccessTokenNotProvidedError)
	})

	t.Run("admin.BanUser: invalid token", func(t *testing.T) {
		_, err := client.BanUser(invalidCtx, &usersv1.BanUserRequest{
			UserId: martin.Id,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthInvalidTokenError)
	})

	t.Run("admin.BanUser: permission denied", func(t *testing.T) {
		_, err := client.BanUser(johnCtx, &usersv1.BanUserRequest{
			UserId: martin.Id,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthPermissionDeniedError)
	})

	t.Run("admin.BanUser: not found", func(t *testing.T) {
		testData := "test_id"
		_, err := client.BanUser(johnAdminCtx, &usersv1.BanUserRequest{
			UserId: testData,
		})

		require.Error(t, err)
		testerror.RequireNotFoundError(t, err, "user", "id", testData)
	})

	t.Run("admin.BanUser: successful", func(t *testing.T) {
		_, err := client.BanUser(johnAdminCtx, &usersv1.BanUserRequest{
			UserId: martin.Id,
		})

		require.NoError(t, err)
	})

	t.Run("admin.GetUserByIdentifier: banned: successful", func(t *testing.T) {
		res, err := client.GetUserByIdentifier(johnAdminCtx, &usersv1.GetUserByIdentifierRequest{
			Identifier: &usersv1.GetUserByIdentifierRequest_UserId{
				UserId: martin.Id,
			},
		})

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, martin.Id, res.Id)
		require.True(t, res.DeletedAt.IsValid())
	})

	t.Run("admin.UnbanUser: access token not provided", func(t *testing.T) {
		_, err := client.UnbanUser(emptyCtx, &usersv1.UnbanUserRequest{
			UserId: martin.Id,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthAccessTokenNotProvidedError)
	})

	t.Run("admin.UnbanUser: invalid token", func(t *testing.T) {
		_, err := client.UnbanUser(invalidCtx, &usersv1.UnbanUserRequest{
			UserId: martin.Id,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthInvalidTokenError)
	})

	t.Run("admin.UnbanUser: permission denied", func(t *testing.T) {
		_, err := client.UnbanUser(johnCtx, &usersv1.UnbanUserRequest{
			UserId: martin.Id,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthPermissionDeniedError)
	})

	t.Run("admin.UnbanUser: not found", func(t *testing.T) {
		testData := "test_id"
		_, err := client.UnbanUser(johnAdminCtx, &usersv1.UnbanUserRequest{
			UserId: testData,
		})

		require.Error(t, err)
		testerror.RequireNotFoundError(t, err, "user", "id", testData)
	})

	t.Run("admin.UnbanUser: successful", func(t *testing.T) {
		_, err := client.UnbanUser(johnAdminCtx, &usersv1.UnbanUserRequest{
			UserId: martin.Id,
		})

		require.NoError(t, err)
	})

	t.Run("admin.GetUserByIdentifier: unbanned: successful", func(t *testing.T) {
		res, err := client.GetUserByIdentifier(johnAdminCtx, &usersv1.GetUserByIdentifierRequest{
			Identifier: &usersv1.GetUserByIdentifierRequest_UserId{
				UserId: martin.Id,
			},
		})

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, martin.Id, res.Id)
		require.False(t, res.DeletedAt.IsValid())
	})

	t.Run("admin.SearchUsers: access token not provided", func(t *testing.T) {
		res, err := client.SearchUsers(emptyCtx, &usersv1.SearchUsersRequest{})

		require.Error(t, err)
		require.Nil(t, res)
		testerror.RequireForbiddenError(t, err, jw.AuthAccessTokenNotProvidedError)
	})

	t.Run("admin.SearchUsers: invalid token", func(t *testing.T) {
		res, err := client.SearchUsers(invalidCtx, &usersv1.SearchUsersRequest{})

		require.Error(t, err)
		require.Nil(t, res)
		testerror.RequireForbiddenError(t, err, jw.AuthInvalidTokenError)
	})

	t.Run("admin.SearchUsers: permission denied", func(t *testing.T) {
		res, err := client.SearchUsers(johnCtx, &usersv1.SearchUsersRequest{})

		require.Error(t, err)
		require.Nil(t, res)
		testerror.RequireForbiddenError(t, err, jw.AuthPermissionDeniedError)
	})

	t.Run("admin.SearchUsers: order by username: sort by desc", func(t *testing.T) {
		order := usersv1.Order_ORDER_USERNAME
		sort := usersv1.Sort_SORT_DESC
		res, err := client.SearchUsers(johnAdminCtx, &usersv1.SearchUsersRequest{
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
		order := usersv1.Order_ORDER_USERNAME
		sort := usersv1.Sort_SORT_ASC
		res, err := client.SearchUsers(johnAdminCtx, &usersv1.SearchUsersRequest{
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
		order := usersv1.Order_ORDER_USERNAME
		sort := usersv1.Sort_SORT_DESC
		res, err := client.SearchUsers(johnAdminCtx, &usersv1.SearchUsersRequest{
			Page:  1,
			Size:  10,
			Order: &order,
			Sort:  &sort,
			UserRating: &usersv1.RatingFiler{
				From: 100,
				To:   100,
			},
		})

		require.Nil(t, err)
		require.Equal(t, 0, len(res.Users))
	})

	t.Run("admin.SearchUsers: filter by rating: successful", func(t *testing.T) {
		order := usersv1.Order_ORDER_USERNAME
		sort := usersv1.Sort_SORT_DESC
		res, err := client.SearchUsers(johnAdminCtx, &usersv1.SearchUsersRequest{
			Page:  1,
			Size:  10,
			Order: &order,
			Sort:  &sort,
			UserRating: &usersv1.RatingFiler{
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
		order := usersv1.Order_ORDER_USERNAME
		sort := usersv1.Sort_SORT_DESC
		res, err := client.SearchUsers(johnAdminCtx, &usersv1.SearchUsersRequest{
			Page:  1,
			Size:  10,
			Order: &order,
			Sort:  &sort,
			UserCoins: &usersv1.CoinsFiler{
				From: 100,
				To:   100,
			},
		})

		require.Nil(t, err)
		require.Equal(t, 0, len(res.Users))
	})

	t.Run("admin.SearchUsers: filter by coins: successful", func(t *testing.T) {
		order := usersv1.Order_ORDER_USERNAME
		sort := usersv1.Sort_SORT_DESC
		res, err := client.SearchUsers(johnAdminCtx, &usersv1.SearchUsersRequest{
			Page:  page,
			Size:  size,
			Order: &order,
			Sort:  &sort,
			UserCoins: &usersv1.CoinsFiler{
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
