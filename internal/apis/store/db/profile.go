package db

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/QuizWars-Ecosystem/go-common/pkg/dbx"
	apperrors "github.com/QuizWars-Ecosystem/go-common/pkg/error"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/profile"
)

func (db *Database) GetProfile(ctx context.Context, userID string) (*profile.Profile, error) {
	builder := dbx.StatementBuilder.
		Select("u.id", "u.username", "u.email", "u.avatar_id", "s.rating", "s.coins", "u.created_at", "u.last_login_at").
		From("users u").
		Join("stats s ON s.user_id = u.id").
		Where(squirrel.Eq{"u.id": userID}).
		Where(squirrel.Eq{"u.deleted_at": nil})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	prof := profile.Profile{
		User: &profile.User{},
	}

	err = db.pool.QueryRow(ctx, query, args...).
		Scan(
			&prof.User.ID,
			&prof.User.Username,
			&prof.Email,
			&prof.User.AvatarID,
			&prof.User.Rating,
			&prof.Coins,
			&prof.User.CreatedAt,
			&prof.User.LastLoginAt,
		)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("user", "id", userID)
	case err != nil:
		return nil, apperrors.Internal(err)
	}

	return &prof, nil
}

func (db *Database) GetUserByID(ctx context.Context, userID string) (*profile.User, error) {
	builder := dbx.StatementBuilder.
		Select("u.id", "u.username", "u.avatar_id", "s.rating", "u.created_at", "u.last_login_at").
		From("users u").
		Join("stats s ON s.user_id = u.id").
		Where(squirrel.Eq{"u.id": userID}).
		Where(squirrel.Eq{"u.deleted_at": nil})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	var user profile.User

	err = db.pool.QueryRow(ctx, query, args...).
		Scan(
			&user.ID,
			&user.Username,
			&user.AvatarID,
			&user.Rating,
			&user.CreatedAt,
			&user.LastLoginAt,
		)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("user", "id", userID)
	case err != nil:
		return nil, apperrors.Internal(err)
	}

	return &user, nil
}

func (db *Database) GetUserByUsername(ctx context.Context, username string) (*profile.User, error) {
	builder := dbx.StatementBuilder.
		Select("u.id", "u.username", "u.avatar_id", "s.rating", "u.created_at", "u.last_login_at").
		From("users u").
		Join("stats s ON s.user_id = u.id").
		Where(squirrel.Eq{"u.username": username}).
		Where(squirrel.Eq{"u.deleted_at": nil})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	var user profile.User

	err = db.pool.QueryRow(ctx, query, args...).
		Scan(
			&user.ID,
			&user.Username,
			&user.AvatarID,
			&user.Rating,
			&user.CreatedAt,
			&user.LastLoginAt,
		)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("user", "username", username)
	case err != nil:
		return nil, apperrors.Internal(err)
	}

	return &user, nil
}

func (db *Database) UpdateProfile(ctx context.Context, userID string, request *profile.UpdateProfile) error {
	builder := dbx.StatementBuilder.
		Update("users").
		Where(squirrel.Eq{"id": userID}).
		Where(squirrel.Eq{"deleted_at": nil})

	if request.Username != nil {
		builder = builder.Set("username", *request.Username)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return apperrors.Internal(err)
	}

	cmd, err := db.pool.Exec(ctx, query, args...)

	switch {
	case dbx.IsUniqueViolation(err, "username"):
		return apperrors.BadRequestHidden(err, "username is already taken")
	case cmd.RowsAffected() == 0:
		return apperrors.NotFound("user", "id", userID)
	case err != nil:
		return apperrors.Internal(err)
	}

	return nil
}

func (db *Database) UpdateProfileAvatar(ctx context.Context, userID string, avatarID int32) error {
	builder := dbx.StatementBuilder.
		Update("users").
		Set("avatar_id", avatarID).
		Where(squirrel.Eq{"id": userID}).
		Where(squirrel.Eq{"deleted_at": nil})

	query, args, err := builder.ToSql()
	if err != nil {
		return apperrors.Internal(err)
	}

	cmd, err := db.pool.Exec(ctx, query, args...)

	switch {
	case err != nil:
		return apperrors.Internal(err)
	case cmd.RowsAffected() == 0:
		return apperrors.NotFound("user", "id", userID)
	}

	return nil
}

func (db *Database) UpdateProfilePassword(ctx context.Context, userID string, password string) error {
	builder := dbx.StatementBuilder.
		Update("users").
		Set("pass_hash", password).
		Where(squirrel.Eq{"id": userID}).
		Where(squirrel.Eq{"deleted_at": nil})

	query, args, err := builder.ToSql()
	if err != nil {
		return apperrors.Internal(err)
	}

	cmd, err := db.pool.Exec(ctx, query, args...)

	switch {
	case err != nil:
		return apperrors.Internal(err)
	case cmd.RowsAffected() == 0:
		return apperrors.NotFound("user", "id", userID)
	}

	return nil
}

func (db *Database) SetProfileRating(ctx context.Context, userID string, rating int32) error {
	builder := dbx.StatementBuilder.
		Update("stats").
		Set("rating", rating).
		Where(squirrel.Eq{"user_id": userID})

	query, args, err := builder.ToSql()
	if err != nil {
		return apperrors.Internal(err)
	}

	cmd, err := db.pool.Exec(ctx, query, args...)
	switch {
	case err != nil:
		return apperrors.Internal(err)
	case cmd.RowsAffected() == 0:
		return apperrors.NotFound("user", "id", userID)
	}

	return nil
}

func (db *Database) SetProfileCoins(ctx context.Context, userID string, coins int64) error {
	builder := dbx.StatementBuilder.
		Update("stats").
		Set("coins", coins).
		Where(squirrel.Eq{"user_id": userID})

	query, args, err := builder.ToSql()
	if err != nil {
		return apperrors.Internal(err)
	}

	cmd, err := db.pool.Exec(ctx, query, args...)
	switch {
	case err != nil:
		return apperrors.Internal(err)
	case cmd.RowsAffected() == 0:
		return apperrors.NotFound("user", "id", userID)
	}

	return nil
}

func (db *Database) DeleteProfile(ctx context.Context, userID string) error {
	builder := dbx.StatementBuilder.
		Update("users").
		Set("deleted_at", time.Now()).
		Where(squirrel.Eq{"id": userID})

	query, args, err := builder.ToSql()
	if err != nil {
		return apperrors.Internal(err)
	}

	cmd, err := db.pool.Exec(ctx, query, args...)

	switch {
	case err != nil:
		return apperrors.Internal(err)
	case cmd.RowsAffected() == 0:
		return apperrors.NotFound("user", "id", userID)
	}

	return nil
}
