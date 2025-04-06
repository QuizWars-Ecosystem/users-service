package modules

import (
	jw "github.com/QuizWars-Ecosystem/go-common/pkg/jwt"
	testerror "github.com/QuizWars-Ecosystem/go-common/pkg/testing/errors"
	"github.com/stretchr/testify/require"
	"testing"

	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
	"github.com/QuizWars-Ecosystem/users-service/tests/integration_tests/config"
)

func SocialServiceTest(t *testing.T, client userspb.UsersSocialServiceClient, _ *config.TestConfig) {
	ctx := t.Context()

	t.Run("social.AddFriend: not found", func(t *testing.T) {
		testData := "test_id"
		_, err := client.AddFriend(ctx, &userspb.AddFriendRequest{
			RequesterId: testData,
			RecipientId: martin.Id,
		})

		require.Error(t, err)
		testerror.RequireNotFoundError(t, err, "user", "id", testData)
	})

	t.Run("social.AddFriend: successful", func(t *testing.T) {
		_, err := client.AddFriend(ctx, &userspb.AddFriendRequest{
			RequesterId: john.Id,
			RecipientId: martin.Id,
		})

		require.NoError(t, err)
	})

	t.Run("social.RejectFriend: not found", func(t *testing.T) {
		_, err := client.RejectFriend(ctx, &userspb.RejectFriendRequest{
			RecipientId: lukas.Id,
			RequesterId: martin.Id,
		})

		require.Error(t, err)
		testerror.RequireNotFoundError(t, err, "user", "id", martin.Id)
	})

	t.Run("social.RejectFriend: successful", func(t *testing.T) {
		_, err := client.RejectFriend(ctx, &userspb.RejectFriendRequest{
			RecipientId: martin.Id,
			RequesterId: john.Id,
		})

		require.NoError(t, err)
	})

	t.Run("social.AddFriend: by list: successful", func(t *testing.T) {
		reqs := []*userspb.AddFriendRequest{
			{
				RequesterId: john.Id,
				RecipientId: martin.Id,
			},
			{
				RequesterId: john.Id,
				RecipientId: lukas.Id,
			},
		}

		for _, req := range reqs {
			_, err := client.AddFriend(ctx, req)

			require.NoError(t, err)
		}
	})

	t.Run("social.AcceptFriend: by list: successful", func(t *testing.T) {
		reqs := []*userspb.AcceptFriendRequest{
			{
				RecipientId: martin.Id,
				RequesterId: john.Id,
			},
			{
				RecipientId: lukas.Id,
				RequesterId: john.Id,
			},
		}

		for _, req := range reqs {
			_, err := client.AcceptFriend(ctx, req)

			require.NoError(t, err)
		}
	})

	t.Run("social.ListFriends: not found", func(t *testing.T) {
		res, err := client.ListFriends(ctx, &userspb.ListFriendsRequest{
			Token:  soniaToken,
			UserId: sonia.Id,
		})

		require.Error(t, err)
		require.Nil(t, res)
		testerror.RequireNotFoundError(t, err, "friends", "id", sonia.Id)
	})

	t.Run("social.ListFriends: by requester id: successful", func(t *testing.T) {
		res, err := client.ListFriends(ctx, &userspb.ListFriendsRequest{
			Token:  johnToken,
			UserId: john.Id,
		})

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Len(t, res.Friends, 2)
	})

	t.Run("social.ListFriends: by recipient id: successful", func(t *testing.T) {
		res, err := client.ListFriends(ctx, &userspb.ListFriendsRequest{
			Token:  martinToken,
			UserId: martin.Id,
		})

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Len(t, res.Friends, 1)
	})

	t.Run("social.BlockFriend: token not provided", func(t *testing.T) {
		_, err := client.BlockFriend(ctx, &userspb.BlockFriendRequest{
			Token:    "",
			UserId:   martin.Id,
			FriendId: john.Id,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthAccessTokenNotProvidedError)
	})

	t.Run("social.BlockFriend: invalid token", func(t *testing.T) {
		_, err := client.BlockFriend(ctx, &userspb.BlockFriendRequest{
			Token:    "invalid token",
			UserId:   martin.Id,
			FriendId: john.Id,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthInvalidTokenError)
	})

	t.Run("social.BlockFriend: permission denied", func(t *testing.T) {
		_, err := client.BlockFriend(ctx, &userspb.BlockFriendRequest{
			Token:    lukasToken,
			UserId:   martin.Id,
			FriendId: john.Id,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthPermissionDeniedError)
	})

	t.Run("social.BlockFriend: not found", func(t *testing.T) {
		testData := "test_id"
		_, err := client.BlockFriend(ctx, &userspb.BlockFriendRequest{
			Token:    martinToken,
			UserId:   martin.Id,
			FriendId: testData,
		})

		require.Error(t, err)
		testerror.RequireNotFoundError(t, err, "user", "id", testData)
	})

	t.Run("social.BlockFriend: successful", func(t *testing.T) {
		_, err := client.BlockFriend(ctx, &userspb.BlockFriendRequest{
			Token:    martinToken,
			UserId:   martin.Id,
			FriendId: john.Id,
		})

		require.NoError(t, err)
	})

	t.Run("social.UnblockFriend: token not provided", func(t *testing.T) {
		_, err := client.UnblockFriend(ctx, &userspb.UnblockFriendRequest{
			Token:    "",
			UserId:   martin.Id,
			FriendId: john.Id,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthAccessTokenNotProvidedError)
	})

	t.Run("social.UnblockFriend: invalid token", func(t *testing.T) {
		_, err := client.BlockFriend(ctx, &userspb.BlockFriendRequest{
			Token:    "invalid token",
			UserId:   martin.Id,
			FriendId: john.Id,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthInvalidTokenError)
	})

	t.Run("social.UnblockFriend: permission denied", func(t *testing.T) {
		_, err := client.BlockFriend(ctx, &userspb.BlockFriendRequest{
			Token:    lukasToken,
			UserId:   martin.Id,
			FriendId: john.Id,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthPermissionDeniedError)
	})

	t.Run("social.UnblockFriend: not found", func(t *testing.T) {
		testData := "test_id"
		_, err := client.BlockFriend(ctx, &userspb.BlockFriendRequest{
			Token:    martinToken,
			UserId:   martin.Id,
			FriendId: testData,
		})

		require.Error(t, err)
		testerror.RequireNotFoundError(t, err, "user", "id", testData)
	})

	t.Run("social.UnblockFriend: successful", func(t *testing.T) {
		_, err := client.BlockFriend(ctx, &userspb.BlockFriendRequest{
			Token:    martinToken,
			UserId:   martin.Id,
			FriendId: john.Id,
		})

		require.NoError(t, err)
	})

	t.Run("social.RemoveFriend: token not provided", func(t *testing.T) {
		_, err := client.RemoveFriend(ctx, &userspb.RemoveFriendRequest{
			Token:       "",
			RequesterId: martin.Id,
			FriendId:    john.Id,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthAccessTokenNotProvidedError)
	})

	t.Run("social.RemoveFriend: invalid token", func(t *testing.T) {
		_, err := client.RemoveFriend(ctx, &userspb.RemoveFriendRequest{
			Token:       "invalid token",
			RequesterId: martin.Id,
			FriendId:    john.Id,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthInvalidTokenError)
	})

	t.Run("social.RemoveFriend: permission denied", func(t *testing.T) {
		_, err := client.RemoveFriend(ctx, &userspb.RemoveFriendRequest{
			Token:       lukasToken,
			RequesterId: martin.Id,
			FriendId:    john.Id,
		})

		require.Error(t, err)
		testerror.RequireForbiddenError(t, err, jw.AuthPermissionDeniedError)
	})

	t.Run("social.RemoveFriend: not found", func(t *testing.T) {
		testData := "test_id"
		_, err := client.RemoveFriend(ctx, &userspb.RemoveFriendRequest{
			Token:       martinToken,
			RequesterId: martin.Id,
			FriendId:    testData,
		})

		require.Error(t, err)
		testerror.RequireNotFoundError(t, err, "user", "id", testData)
	})

	t.Run("social.RemoveFriend: successful", func(t *testing.T) {
		_, err := client.RemoveFriend(ctx, &userspb.RemoveFriendRequest{
			Token:       johnToken,
			RequesterId: john.Id,
			FriendId:    martin.Id,
		})

		require.NoError(t, err)
	})

	t.Run("social.ListFriends: by requester id: successful", func(t *testing.T) {
		res, err := client.ListFriends(ctx, &userspb.ListFriendsRequest{
			Token:  johnToken,
			UserId: john.Id,
		})

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Len(t, res.Friends, 1)
	})
}
