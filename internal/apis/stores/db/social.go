package db

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/QuizWars-Ecosystem/go-common/pkg/dbx"
	apperrors "github.com/QuizWars-Ecosystem/go-common/pkg/error"
	"github.com/QuizWars-Ecosystem/users-service/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Social struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewSocial(db *pgxpool.Pool, logger *zap.Logger) *Social {
	return &Social{
		db:     db,
		logger: logger,
	}
}

func (s *Social) AddFriend(ctx context.Context, userId string, friendId string) error {
	b := &pgx.Batch{}

	builder := dbx.StatementBuilder.
		Insert("friends").
		Columns("user_id", "friend_id").
		Values(userId, friendId).
		Suffix("ON CONFLICT DO NOTHING")

	if err := dbx.QueryBatch(b, builder); err != nil {
		return apperrors.Internal(err)
	}

	builder = dbx.StatementBuilder.
		Insert("friends").
		Columns("user_id", "friend_id").
		Values(friendId, userId).
		Suffix("ON CONFLICT DO NOTHING")

	if err := dbx.QueryBatch(b, builder); err != nil {
		return apperrors.Internal(err)
	}

	err := s.db.SendBatch(ctx, b).Close()

	switch {
	case dbx.IsNoRows(err):
		return apperrors.NotFound("user", "id", friendId)
	case err != nil:
		return apperrors.Internal(err)
	}

	return nil
}

func (s *Social) AcceptFriend(ctx context.Context, userId string, friendId string) error {
	b := &pgx.Batch{}

	builder := dbx.StatementBuilder.
		Update("friends").
		Set("status", models.Accepted).
		Where(squirrel.Eq{"user_id": userId, "friend_id": friendId})

	if err := dbx.QueryBatch(b, builder); err != nil {
		return apperrors.Internal(err)
	}

	builder = dbx.StatementBuilder.
		Update("friends").
		Set("status", models.Accepted).
		Where(squirrel.Eq{"user_id": friendId, "friend_id": userId})

	if err := dbx.QueryBatch(b, builder); err != nil {
		return apperrors.Internal(err)
	}

	err := s.db.SendBatch(ctx, b).Close()
	if err != nil {
		return apperrors.Internal(err)
	}

	return nil

}

func (s *Social) RemoveFriend(ctx context.Context, userId string, friendId string) error {
	b := &pgx.Batch{}

	builder := dbx.StatementBuilder.
		Delete("friends").
		Where(squirrel.Eq{"user_id": userId, "friend_id": friendId})

	if err := dbx.QueryBatch(b, builder); err != nil {
		return apperrors.Internal(err)
	}

	builder = dbx.StatementBuilder.
		Delete("friends").
		Where(squirrel.Eq{"user_id": friendId, "friend_id": userId})

	if err := dbx.QueryBatch(b, builder); err != nil {
		return apperrors.Internal(err)
	}

	err := s.db.SendBatch(ctx, b).Close()
	switch {
	case dbx.IsNoRows(err):
		return nil
	case err != nil:
		return apperrors.Internal(err)
	}

	return nil
}

func (s *Social) BanFriend(ctx context.Context, userId string, friendId string) error {
	b := &pgx.Batch{}

	builder := dbx.StatementBuilder.
		Update("friends").
		Set("status", models.Blocked).
		Where(squirrel.Eq{"user_id": userId, "friend_id": friendId})

	if err := dbx.QueryBatch(b, builder); err != nil {
		return apperrors.Internal(err)
	}

	builder = dbx.StatementBuilder.
		Update("friends").
		Set("status", models.Blocked).
		Where(squirrel.Eq{"user_id": friendId, "friend_id": userId})

	if err := dbx.QueryBatch(b, builder); err != nil {
		return apperrors.Internal(err)
	}

	err := s.db.SendBatch(ctx, b).Close()
	if err != nil {
		return apperrors.Internal(err)
	}

	return nil

}

func (s *Social) UnbanFriend(ctx context.Context, userId string, friendId string) error {
	b := &pgx.Batch{}

	builder := dbx.StatementBuilder.
		Update("friends").
		Set("status", models.Accepted).
		Where(squirrel.Eq{"user_id": userId, "friend_id": friendId})

	if err := dbx.QueryBatch(b, builder); err != nil {
		return apperrors.Internal(err)
	}

	builder = dbx.StatementBuilder.
		Update("friends").
		Set("status", models.Accepted).
		Where(squirrel.Eq{"user_id": friendId, "friend_id": userId})

	if err := dbx.QueryBatch(b, builder); err != nil {
		return apperrors.Internal(err)
	}

	err := s.db.SendBatch(ctx, b).Close()
	if err != nil {
		return apperrors.Internal(err)
	}

	return nil

}
