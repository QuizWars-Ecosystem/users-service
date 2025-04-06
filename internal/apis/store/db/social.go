package db

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/QuizWars-Ecosystem/go-common/pkg/dbx"
	apperrors "github.com/QuizWars-Ecosystem/go-common/pkg/error"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/profile"
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

func (s *Social) AddFriend(ctx context.Context, requesterID, recipientID string) error {
	builder := dbx.StatementBuilder.
		Insert("friends").
		Columns("user_id", "friend_id").
		Values(requesterID, recipientID).
		Suffix("ON CONFLICT DO NOTHING")

	query, args, err := builder.ToSql()
	if err != nil {
		return apperrors.Internal(err)
	}

	_, err = s.db.Exec(ctx, query, args...)

	switch {
	case dbx.IsForeignKeyViolation(err, "user_id"):
		return apperrors.NotFound("user", "id", requesterID)
	case dbx.IsForeignKeyViolation(err, "friend_id"):
		return apperrors.NotFound("user", "id", requesterID)
	case err != nil:
		return apperrors.Internal(err)
	}

	return nil
}

func (s *Social) AcceptFriend(ctx context.Context, recipientID, requesterID string) error {
	builder := dbx.StatementBuilder.
		Update("friends").
		Set("status", "accepted").
		Where(squirrel.Eq{"friend_id": recipientID}).
		Where(squirrel.Eq{"user_id": requesterID})

	query, args, err := builder.ToSql()
	if err != nil {
		return apperrors.Internal(err)
	}

	_, err = s.db.Exec(ctx, query, args...)
	if err != nil {
		return apperrors.Internal(err)
	}

	return nil
}

func (s *Social) RejectFriend(ctx context.Context, recipientID, requesterID string) error {
	builder := dbx.StatementBuilder.
		Delete("friends").
		Where(squirrel.Eq{"friend_id": recipientID}).
		Where(squirrel.Eq{"user_id": requesterID})

	query, args, err := builder.ToSql()
	if err != nil {
		return apperrors.Internal(err)
	}

	cmd, err := s.db.Exec(ctx, query, args...)
	switch {
	case dbx.IsForeignKeyViolation(err, "user_id"):
		return apperrors.NotFound("user", "id", recipientID)
	case dbx.IsForeignKeyViolation(err, "friend_id"):
		return apperrors.NotFound("user", "id", requesterID)
	case cmd.RowsAffected() == 0:
		return apperrors.NotFound("user", "id", requesterID)
	case err != nil:
		return apperrors.Internal(err)
	}

	return nil
}

func (s *Social) RemoveFriend(ctx context.Context, userID string, friendID string) error {
	builder := dbx.StatementBuilder.
		Delete("friends").
		Where(
			squirrel.Or{
				squirrel.And{squirrel.Eq{"user_id": userID}, squirrel.Eq{"friend_id": friendID}},
				squirrel.And{squirrel.Eq{"user_id": friendID}, squirrel.Eq{"friend_id": userID}},
			},
		)

	query, args, err := builder.ToSql()
	if err != nil {
		return apperrors.Internal(err)
	}

	cmd, err := s.db.Exec(ctx, query, args...)
	switch {
	case dbx.IsForeignKeyViolation(err, "user_id"):
		return apperrors.NotFound("user", "id", userID)
	case dbx.IsForeignKeyViolation(err, "friend_id"):
		return apperrors.NotFound("user", "id", friendID)
	case cmd.RowsAffected() == 0:
		return apperrors.NotFound("user", "id", friendID)
	case err != nil:
		return apperrors.Internal(err)
	}

	return nil
}

func (s *Social) GetFriends(ctx context.Context, userID string) ([]*profile.Friend, error) {
	builder := dbx.StatementBuilder.
		Select("u.id", "u.avatar_id", "u.username", "s.rating", "u.created_at", "u.last_login_at", "f.status").
		From("users u").
		JoinClause("JOIN friends f ON (f.user_id = ? AND u.id = f.friend_id) OR (f.friend_id = ? AND u.id = f.user_id)", userID, userID).
		Join("stats s ON s.user_id = u.id").
		Where(squirrel.Eq{"u.deleted_at": nil}).
		Where(squirrel.NotEq{"u.id": userID}).
		Where(squirrel.Or{
			squirrel.NotEq{"f.status": "pending"},
			squirrel.And{
				squirrel.NotEq{"f.user_id": userID},
				squirrel.Eq{"f.status": "pending"},
			},
		}).
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

	switch {
	case rows.Err() != nil:
		return nil, apperrors.Internal(rows.Err())
	case len(friends) == 0:
		return nil, apperrors.NotFound("friends", "id", userID)
	}

	return friends, nil
}

func (s *Social) BanFriend(ctx context.Context, userID string, friendID string) error {
	builder := dbx.StatementBuilder.
		Update("friends").
		Set("status", "blocked").
		Where(
			squirrel.Or{
				squirrel.And{squirrel.Eq{"user_id": userID}, squirrel.Eq{"friend_id": friendID}},
				squirrel.And{squirrel.Eq{"user_id": friendID}, squirrel.Eq{"friend_id": userID}},
			},
		)

	query, args, err := builder.ToSql()
	if err != nil {
		return apperrors.Internal(err)
	}

	cmd, err := s.db.Exec(ctx, query, args...)
	switch {
	case dbx.IsForeignKeyViolation(err, "user_id"):
		return apperrors.NotFound("user", "id", userID)
	case dbx.IsForeignKeyViolation(err, "friend_id"):
		return apperrors.NotFound("user", "id", friendID)
	case cmd.RowsAffected() == 0:
		return apperrors.NotFound("user", "id", friendID)
	case err != nil:
		return apperrors.Internal(err)
	}

	return nil
}

func (s *Social) UnbanFriend(ctx context.Context, userID string, friendID string) error {
	builder := dbx.StatementBuilder.
		Update("friends").
		Set("status", "accepted").
		Where(
			squirrel.Or{
				squirrel.And{squirrel.Eq{"user_id": userID}, squirrel.Eq{"friend_id": friendID}},
				squirrel.And{squirrel.Eq{"user_id": friendID}, squirrel.Eq{"friend_id": userID}},
			},
		)

	query, args, err := builder.ToSql()
	if err != nil {
		return apperrors.Internal(err)
	}

	cmd, err := s.db.Exec(ctx, query, args...)
	switch {
	case dbx.IsForeignKeyViolation(err, "user_id"):
		return apperrors.NotFound("user", "id", userID)
	case dbx.IsForeignKeyViolation(err, "friend_id"):
		return apperrors.NotFound("user", "id", friendID)
	case cmd.RowsAffected() == 0:
		return apperrors.NotFound("user", "id", friendID)
	case err != nil:
		return apperrors.Internal(err)
	}

	return nil
}
