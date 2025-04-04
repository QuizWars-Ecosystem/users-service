package db

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/QuizWars-Ecosystem/go-common/pkg/dbx"
	apperrors "github.com/QuizWars-Ecosystem/go-common/pkg/error"
	"github.com/QuizWars-Ecosystem/users-service/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

type Profile struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewProfile(db *pgxpool.Pool, logger *zap.Logger) *Profile {
	return &Profile{
		db:     db,
		logger: logger,
	}
}

func (p *Profile) GetProfile(ctx context.Context, userID string) (*models.Profile, error) {
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

	var profile models.Profile

	err = p.db.QueryRow(ctx, query, args...).
		Scan(
			&profile.User.ID,
			&profile.User.Username,
			&profile.Email,
			&profile.User.AvatarID,
			&profile.User.Rating,
			&profile.Coins,
			&profile.User.CreatedAt,
			&profile.User.LastLoginAt,
		)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("user", "id", userID)
	case err != nil:
		return nil, apperrors.Internal(err)
	}

	return &profile, nil
}

func (p *Profile) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
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

	var user models.User

	err = p.db.QueryRow(ctx, query, args...).
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

func (p *Profile) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
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

	var user models.User

	err = p.db.QueryRow(ctx, query, args...).
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

func (p *Profile) UpdateProfile(ctx context.Context, userID string, request *models.UpdateProfile) error {
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

	cmd, err := p.db.Exec(ctx, query, args...)

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

func (p *Profile) UpdateProfileAvatar(ctx context.Context, userID string, avatarID int32) error {
	builder := dbx.StatementBuilder.
		Update("users").
		Set("avatar_id", avatarID).
		Where(squirrel.Eq{"id": userID}).
		Where(squirrel.Eq{"deleted_at": nil})

	query, args, err := builder.ToSql()
	if err != nil {
		return apperrors.Internal(err)
	}

	cmd, err := p.db.Exec(ctx, query, args...)

	switch {
	case err != nil:
		return apperrors.Internal(err)
	case cmd.RowsAffected() == 0:
		return apperrors.NotFound("user", "id", userID)
	}

	return nil
}

func (p *Profile) UpdateProfilePassword(ctx context.Context, userID string, password string) error {
	builder := dbx.StatementBuilder.
		Update("users").
		Set("pass_hash", password).
		Where(squirrel.Eq{"id": userID}).
		Where(squirrel.Eq{"deleted_at": nil})

	query, args, err := builder.ToSql()
	if err != nil {
		return apperrors.Internal(err)
	}

	cmd, err := p.db.Exec(ctx, query, args...)

	switch {
	case err != nil:
		return apperrors.Internal(err)
	case cmd.RowsAffected() == 0:
		return apperrors.NotFound("user", "id", userID)
	}

	return nil
}

func (p *Profile) SetProfileRating(ctx context.Context, userID string, rating int32) error {
	builder := dbx.StatementBuilder.
		Update("stats").
		Set("rating", rating).
		Where(squirrel.Eq{"user_id": userID})

	query, args, err := builder.ToSql()
	if err != nil {
		return apperrors.Internal(err)
	}

	cmd, err := p.db.Exec(ctx, query, args...)
	switch {
	case err != nil:
		return apperrors.Internal(err)
	case cmd.RowsAffected() == 0:
		return apperrors.NotFound("user", "id", userID)
	}

	return nil
}

func (p *Profile) SetProfileCoins(ctx context.Context, userID string, coins int64) error {
	builder := dbx.StatementBuilder.
		Update("stats").
		Set("coins", coins).
		Where(squirrel.Eq{"user_id": userID})

	query, args, err := builder.ToSql()
	if err != nil {
		return apperrors.Internal(err)
	}

	cmd, err := p.db.Exec(ctx, query, args...)
	switch {
	case err != nil:
		return apperrors.Internal(err)
	case cmd.RowsAffected() == 0:
		return apperrors.NotFound("user", "id", userID)
	}

	return nil
}

func (p *Profile) DeleteProfile(ctx context.Context, userID string) error {
	builder := dbx.StatementBuilder.
		Update("users").
		Set("deleted_at", time.Now()).
		Where(squirrel.Eq{"id": userID})

	query, args, err := builder.ToSql()
	if err != nil {
		return apperrors.Internal(err)
	}

	cmd, err := p.db.Exec(ctx, query, args...)

	switch {
	case err != nil:
		return apperrors.Internal(err)
	case cmd.RowsAffected() == 0:
		return apperrors.NotFound("user", "id", userID)
	}

	return nil
}
