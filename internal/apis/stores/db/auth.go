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

type Auth struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewAuth(db *pgxpool.Pool, logger *zap.Logger) *Auth {
	return &Auth{
		db:     db,
		logger: logger,
	}
}

func (a *Auth) SaveProfile(ctx context.Context, profile *models.ProfileWithCredits) (*models.Profile, error) {
	builder := dbx.StatementBuilder.
		Insert("users").
		Columns("username", "email", "pass_hash", "avatar_id").
		Values(profile.Profile.User.Username, profile.Profile.Email, profile.Password, profile.Profile.User.AvatarID).
		Suffix("RETURNING id, created_at")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	err = a.db.QueryRow(ctx, query, args...).
		Scan(profile.Profile.User.ID, profile.Profile.User.CreatedAt)

	if err != nil {
		builder = dbx.StatementBuilder.
			Insert("stats").
			Columns("user_id").
			Values(profile.Profile.User.ID)

		query, args, err = builder.ToSql()
		if err != nil {
			return nil, apperrors.Internal(err)
		}

		_, err = a.db.Exec(ctx, query, args...)
		if err != nil {
			return nil, apperrors.Internal(err)
		}
	}

	switch {
	case dbx.IsUniqueViolation(err, "username"):
		return nil, apperrors.BadRequestHidden(err, "username is already taken")
	case dbx.IsUniqueViolation(err, "email"):
		return nil, apperrors.BadRequestHidden(err, "email is already taken")
	case err != nil:
		return nil, apperrors.Internal(err)
	}

	return profile.Profile, nil
}

func (a *Auth) GetProfileByUsername(ctx context.Context, username string) (*models.ProfileWithCredits, error) {
	builder := dbx.StatementBuilder.
		Select("u.id", "u.username", "u.email", "u.pass_hash", "u.avatar_id", "s.rating", "s.coins", "u.created_at", "u.last_login_at").
		From("users u").
		Join("stats s ON s.user_id = u.id").
		Where(squirrel.Eq{"username": username}).
		Where(squirrel.Eq{"deleted_at": nil})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	profile := models.ProfileWithCredits{
		Profile: &models.Profile{
			User: &models.User{},
		},
	}

	err = a.db.QueryRow(ctx, query, args...).
		Scan(
			&profile.Profile.User.ID,
			&profile.Profile.User.Username,
			&profile.Profile.Email,
			&profile.Password,
			&profile.Profile.User.AvatarID,
			&profile.Profile.User.Rating,
			&profile.Profile.Coins,
			&profile.Profile.User.CreatedAt,
			&profile.Profile.User.LastLoginAt,
		)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("user", "username", username)
	case err != nil:
		return nil, apperrors.Internal(err)
	}

	return &profile, nil
}

func (a *Auth) GetProfileByEmail(ctx context.Context, email string) (*models.ProfileWithCredits, error) {
	builder := dbx.StatementBuilder.
		Select("u.id", "u.username", "u.email", "u.pass_hash", "u.avatar_id", "s.rating", "s.coins", "u.created_at", "u.last_login_at").
		From("users u").
		Join("stats s ON s.user_id = u.id").
		Where(squirrel.Eq{"email": email}).
		Where(squirrel.Eq{"deleted_at": nil})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	profile := models.ProfileWithCredits{
		Profile: &models.Profile{
			User: &models.User{},
		},
	}

	err = a.db.QueryRow(ctx, query, args...).
		Scan(
			&profile.Profile.User.ID,
			&profile.Profile.User.Username,
			&profile.Profile.Email,
			&profile.Password,
			&profile.Profile.User.AvatarID,
			&profile.Profile.User.Rating,
			&profile.Profile.Coins,
			&profile.Profile.User.CreatedAt,
			&profile.Profile.User.LastLoginAt,
		)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("user", "email", email)
	case err != nil:
		return nil, apperrors.Internal(err)
	}

	return &profile, nil
}

func (a *Auth) SetLastLogin(ctx context.Context, userID string) error {
	builder := dbx.StatementBuilder.
		Update("users").
		Set("last_login_at", time.Now()).
		Where(squirrel.Eq{"id": userID}).
		Where(squirrel.Eq{"deleted_at": nil})

	query, args, err := builder.ToSql()
	if err != nil {
		return apperrors.Internal(err)
	}

	_, err = a.db.Exec(ctx, query, args...)
	switch {
	case dbx.IsNoRows(err):
		return apperrors.NotFound("user", "id", userID)
	case err != nil:
		return apperrors.Internal(err)
	}

	return nil
}
