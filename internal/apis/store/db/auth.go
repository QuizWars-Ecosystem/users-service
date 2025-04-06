package db

import (
	"context"
	"time"

	"github.com/QuizWars-Ecosystem/users-service/internal/models/auth"

	"github.com/Masterminds/squirrel"
	"github.com/QuizWars-Ecosystem/go-common/pkg/dbx"
	apperrors "github.com/QuizWars-Ecosystem/go-common/pkg/error"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/profile"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
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

func (a *Auth) SaveProfile(ctx context.Context, p *auth.ProfileWithCredentials) (*profile.Profile, error) {
	builder := dbx.StatementBuilder.
		Insert("users").
		Columns("id", "username", "email", "pass_hash", "avatar_id", "created_at").
		Values(p.Profile.User.ID, p.Profile.User.Username, p.Profile.Email, p.Password, p.Profile.User.AvatarID, time.Now())

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	_, err = a.db.Exec(ctx, query, args...)

	if err == nil {
		builder = dbx.StatementBuilder.
			Insert("stats").
			Columns("user_id").
			Values(p.Profile.User.ID).
			Suffix("ON CONFLICT DO NOTHING")

		query, args, err = builder.ToSql()
		if err != nil {
			return nil, apperrors.Internal(err)
		}

		_, err = a.db.Exec(ctx, query, args...)
		switch {
		case dbx.IsForeignKeyViolation(err, "user_id"):
			return nil, apperrors.BadRequestHidden(err, "user already exists")
		case err != nil:
			return nil, apperrors.Internal(err)
		}
	}

	switch {
	case dbx.IsUniqueViolation(err, "username"):
		return nil, apperrors.AlreadyExists("user", "username", p.Profile.User.Username)
	case dbx.IsUniqueViolation(err, "email"):
		return nil, apperrors.AlreadyExists("user", "email", p.Profile.Email)
	case err != nil:
		return nil, apperrors.Internal(err)
	}

	return p.Profile, nil
}

func (a *Auth) GetProfileByUsername(ctx context.Context, username string) (*auth.ProfileWithCredentials, error) {
	builder := dbx.StatementBuilder.
		Select("u.id", "u.username", "u.email", "u.pass_hash", "u.role", "u.avatar_id", "s.rating", "s.coins", "u.created_at", "u.last_login_at").
		From("users u").
		Join("stats s ON s.user_id = u.id").
		Where(squirrel.Eq{"u.username": username}).
		Where(squirrel.Eq{"u.deleted_at": nil})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	p := auth.ProfileWithCredentials{
		Profile: &profile.Profile{
			User: &profile.User{},
		},
	}

	err = a.db.QueryRow(ctx, query, args...).
		Scan(
			&p.Profile.User.ID,
			&p.Profile.User.Username,
			&p.Profile.Email,
			&p.Password,
			&p.Role,
			&p.Profile.User.AvatarID,
			&p.Profile.User.Rating,
			&p.Profile.Coins,
			&p.Profile.User.CreatedAt,
			&p.Profile.User.LastLoginAt,
		)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("user", "username", username)
	case err != nil:
		return nil, apperrors.Internal(err)
	}

	return &p, nil
}

func (a *Auth) GetProfileByEmail(ctx context.Context, email string) (*auth.ProfileWithCredentials, error) {
	builder := dbx.StatementBuilder.
		Select("u.id", "u.username", "u.email", "u.pass_hash", "u.role", "u.avatar_id", "s.rating", "s.coins", "u.created_at", "u.last_login_at").
		From("users u").
		Join("stats s ON s.user_id = u.id").
		Where(squirrel.Eq{"u.email": email}).
		Where(squirrel.Eq{"u.deleted_at": nil})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	p := auth.ProfileWithCredentials{
		Profile: &profile.Profile{
			User: &profile.User{},
		},
	}

	err = a.db.QueryRow(ctx, query, args...).
		Scan(
			&p.Profile.User.ID,
			&p.Profile.User.Username,
			&p.Profile.Email,
			&p.Password,
			&p.Role,
			&p.Profile.User.AvatarID,
			&p.Profile.User.Rating,
			&p.Profile.Coins,
			&p.Profile.User.CreatedAt,
			&p.Profile.User.LastLoginAt,
		)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("user", "email", email)
	case err != nil:
		return nil, apperrors.Internal(err)
	}

	return &p, nil
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

	cmd, err := a.db.Exec(ctx, query, args...)
	switch {
	case err != nil:
		return apperrors.Internal(err)
	case cmd.RowsAffected() == 0:
		return apperrors.NotFound("user", "id", userID)
	}

	return nil
}
