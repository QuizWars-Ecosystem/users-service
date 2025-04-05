package db

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/QuizWars-Ecosystem/go-common/pkg/dbx"
	apperrors "github.com/QuizWars-Ecosystem/go-common/pkg/error"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/profile"
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

func (s *Social) AddFriend(ctx context.Context, userID string, friendID string) error {
	b := &pgx.Batch{}

	builder := dbx.StatementBuilder.
		Insert("friends").
		Columns("user_id", "friend_id").
		Values(userID, friendID).
		Suffix("ON CONFLICT DO NOTHING")

	if err := dbx.QueryBatch(b, builder); err != nil {
		return apperrors.Internal(err)
	}

	builder = dbx.StatementBuilder.
		Insert("friends").
		Columns("user_id", "friend_id").
		Values(friendID, userID).
		Suffix("ON CONFLICT DO NOTHING")

	if err := dbx.QueryBatch(b, builder); err != nil {
		return apperrors.Internal(err)
	}

	err := s.db.SendBatch(ctx, b).Close()

	switch {
	case dbx.IsNoRows(err):
		return apperrors.NotFound("user", "id", friendID)
	case err != nil:
		return apperrors.Internal(err)
	}

	return nil
}

func (s *Social) AcceptFriend(ctx context.Context, userID string, friendID string) error {
	b := &pgx.Batch{}

	builder := dbx.StatementBuilder.
		Update("friends").
		Set("status", profile.Accepted).
		Where(squirrel.Eq{"user_id": userID, "friend_id": friendID})

	if err := dbx.QueryBatch(b, builder); err != nil {
		return apperrors.Internal(err)
	}

	builder = dbx.StatementBuilder.
		Update("friends").
		Set("status", profile.Accepted).
		Where(squirrel.Eq{"user_id": friendID, "friend_id": userID})

	if err := dbx.QueryBatch(b, builder); err != nil {
		return apperrors.Internal(err)
	}

	err := s.db.SendBatch(ctx, b).Close()
	if err != nil {
		return apperrors.Internal(err)
	}

	return nil
}

func (s *Social) RemoveFriend(ctx context.Context, userID string, friendID string) error {
	b := &pgx.Batch{}

	builder := dbx.StatementBuilder.
		Delete("friends").
		Where(squirrel.Eq{"user_id": userID, "friend_id": friendID})

	if err := dbx.QueryBatch(b, builder); err != nil {
		return apperrors.Internal(err)
	}

	builder = dbx.StatementBuilder.
		Delete("friends").
		Where(squirrel.Eq{"user_id": friendID, "friend_id": userID})

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

func (s *Social) GetFriends(ctx context.Context, userID string) ([]*profile.Friend, error) {
	builder := dbx.StatementBuilder.
		Select("u.id", "u.avatar_id", "u.username", "s.rating", "u.created_at", "u.last_login_at", "f.status").
		From("users u").
		Join("friends f ON (f.user_id = $1 AND u.id = f.friend_id) OR (f.friend_id = $1 AND u.id = f.user_id)", userID).
		Join("stats s ON s.user_id = f.friend_id").
		Where(squirrel.Eq{"u.deleted_at": nil}).
		Where(squirrel.NotEq{"u.id": userID}).
		OrderBy("u.username DESC")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	var friends []*profile.Friend

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	defer rows.Close()

	for rows.Next() {
		f := profile.Friend{
			User: &profile.User{},
		}

		if err = rows.Scan(
			&f.User.ID,
			&f.User.AvatarID,
			&f.User.Username,
			&f.User.Rating,
			&f.User.CreatedAt,
			&f.User.LastLoginAt,
			&f.Status,
		); err != nil {
			return nil, apperrors.Internal(err)
		}

		friends = append(friends, &f)
	}

	if err = rows.Err(); err != nil {
		return nil, apperrors.Internal(err)
	}

	return friends, nil
}

func (s *Social) BanFriend(ctx context.Context, userID string, friendID string) error {
	b := &pgx.Batch{}

	builder := dbx.StatementBuilder.
		Update("friends").
		Set("status", profile.Blocked).
		Where(squirrel.Eq{"user_id": userID, "friend_id": friendID})

	if err := dbx.QueryBatch(b, builder); err != nil {
		return apperrors.Internal(err)
	}

	builder = dbx.StatementBuilder.
		Update("friends").
		Set("status", profile.Blocked).
		Where(squirrel.Eq{"user_id": friendID, "friend_id": userID})

	if err := dbx.QueryBatch(b, builder); err != nil {
		return apperrors.Internal(err)
	}

	err := s.db.SendBatch(ctx, b).Close()
	if err != nil {
		return apperrors.Internal(err)
	}

	return nil
}

func (s *Social) UnbanFriend(ctx context.Context, userID string, friendID string) error {
	b := &pgx.Batch{}

	builder := dbx.StatementBuilder.
		Update("friends").
		Set("status", profile.Accepted).
		Where(squirrel.Eq{"user_id": userID, "friend_id": friendID})

	if err := dbx.QueryBatch(b, builder); err != nil {
		return apperrors.Internal(err)
	}

	builder = dbx.StatementBuilder.
		Update("friends").
		Set("status", profile.Accepted).
		Where(squirrel.Eq{"user_id": friendID, "friend_id": userID})

	if err := dbx.QueryBatch(b, builder); err != nil {
		return apperrors.Internal(err)
	}

	err := s.db.SendBatch(ctx, b).Close()
	if err != nil {
		return apperrors.Internal(err)
	}

	return nil
}
