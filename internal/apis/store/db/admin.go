package db

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/QuizWars-Ecosystem/go-common/pkg/dbx"
	apperrors "github.com/QuizWars-Ecosystem/go-common/pkg/error"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/admin"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/profile"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Admin struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewAdmin(db *pgxpool.Pool, logger *zap.Logger) *Admin {
	return &Admin{
		db:     db,
		logger: logger,
	}
}

func (a *Admin) SearchUsers(ctx context.Context, filter *admin.SearchFilter) ([]*profile.UserAdmin, int, error) {
	builder := dbx.StatementBuilder.
		Select("u.id", "u.username", "u.email", "u.avatar_id", "s.rating", "s.coins", "u.created_at", "u.last_login_at", "u.deleted_at").
		From("users u").
		Join("stats s ON s.user_id = u.id").
		OrderBy(filter.Order.String() + " " + filter.Sort.String()).
		Limit(filter.Limit).
		Offset(filter.Offset)

	if filter.RatingFilter != nil {
		builder = builder.
			Where(squirrel.GtOrEq{"s.rating": filter.RatingFilter.From}).
			Where(squirrel.LtOrEq{"s.rating": filter.RatingFilter.To})
	}

	if filter.CoinsFilter != nil {
		builder = builder.
			Where(squirrel.GtOrEq{"s.coins": filter.CoinsFilter.From}).
			Where(squirrel.LtOrEq{"s.coins": filter.CoinsFilter.To})
	}

	if filter.CreatedAtFilter != nil {
		builder = builder.
			Where(squirrel.GtOrEq{"u.created_at": filter.CreatedAtFilter.From}).
			Where(squirrel.LtOrEq{"u.created_at": filter.CreatedAtFilter.To})
	}

	if filter.DeletedAtFilter != nil {
		builder = builder.
			Where(squirrel.GtOrEq{"u.deleted_at": filter.DeletedAtFilter.From}).
			Where(squirrel.LtOrEq{"u.deleted_at": filter.DeletedAtFilter.To})
	}

	countQuery := dbx.StatementBuilder.Select("COUNT(*)").From("users")

	b := &pgx.Batch{}

	if err := dbx.QueryBatch(b, builder); err != nil {
		return nil, 0, apperrors.Internal(err)
	}

	if err := dbx.QueryBatch(b, countQuery); err != nil {
		return nil, 0, apperrors.Internal(err)
	}

	br := a.db.SendBatch(ctx, b)
	defer func() {
		_ = br.Close()
	}()

	rows, err := br.Query()
	if err != nil {
		return nil, 0, apperrors.Internal(err)
	}

	defer rows.Close()

	var users []*profile.UserAdmin

	for rows.Next() {
		u := profile.UserAdmin{
			Profile: &profile.Profile{
				User: &profile.User{},
			},
		}

		if err = rows.Scan(
			&u.Profile.User.ID,
			&u.Profile.User.Username,
			&u.Profile.Email,
			&u.Profile.User.AvatarID,
			&u.Profile.User.Rating,
			&u.Profile.Coins,
			&u.Profile.User.CreatedAt,
			&u.Profile.User.LastLoginAt,
			&u.DeletedAt,
		); err != nil {
			return nil, 0, apperrors.Internal(err)
		}

		users = append(users, &u)
	}

	if rows.Err() != nil {
		return nil, 0, apperrors.Internal(err)
	}

	var total int
	if err = br.QueryRow().Scan(&total); err != nil {
		return nil, 0, apperrors.Internal(err)
	}

	return users, total, nil
}

func (a *Admin) GetUserByID(ctx context.Context, userID string) (*profile.UserAdmin, error) {
	builder := dbx.StatementBuilder.
		Select("u.id", "u.username", "u.email", "u.avatar_id", "s.rating", "s.coins", "u.created_at", "u.last_login_at", "u.deleted_at").
		From("users u").
		Join("stats s ON s.user_id = u.id").
		Where(squirrel.Eq{"u.id": userID})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	u := profile.UserAdmin{
		Profile: &profile.Profile{
			User: &profile.User{},
		},
	}

	err = a.db.QueryRow(ctx, query, args...).
		Scan(
			&u.Profile.User.ID,
			&u.Profile.User.Username,
			&u.Profile.Email,
			&u.Profile.User.AvatarID,
			&u.Profile.User.Rating,
			&u.Profile.Coins,
			&u.Profile.User.CreatedAt,
			&u.Profile.User.LastLoginAt,
			&u.DeletedAt,
		)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("user", "id", userID)
	case err != nil:
		return nil, apperrors.Internal(err)
	}

	return &u, nil
}

func (a *Admin) GetUserByUsername(ctx context.Context, username string) (*profile.UserAdmin, error) {
	builder := dbx.StatementBuilder.
		Select("u.id", "u.username", "u.email", "u.avatar_id", "s.rating", "s.coins", "u.created_at", "u.last_login_at", "u.deleted_at").
		From("users u").
		Join("stats s ON s.user_id = u.id").
		Where(squirrel.Eq{"u.username": username})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	u := profile.UserAdmin{
		Profile: &profile.Profile{
			User: &profile.User{},
		},
	}

	err = a.db.QueryRow(ctx, query, args...).
		Scan(
			&u.Profile.User.ID,
			&u.Profile.User.Username,
			&u.Profile.Email,
			&u.Profile.User.AvatarID,
			&u.Profile.User.Rating,
			&u.Profile.Coins,
			&u.Profile.User.CreatedAt,
			&u.Profile.User.LastLoginAt,
			&u.DeletedAt,
		)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("user", "username", username)
	case err != nil:
		return nil, apperrors.Internal(err)
	}

	return &u, nil
}

func (a *Admin) GetUserByEmail(ctx context.Context, email string) (*profile.UserAdmin, error) {
	builder := dbx.StatementBuilder.
		Select("u.id", "u.username", "u.email", "u.avatar_id", "s.rating", "s.coins", "u.created_at", "u.last_login_at", "u.deleted_at").
		From("users u").
		Join("stats s ON s.user_id = u.id").
		Where(squirrel.Eq{"u.email": email})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	u := profile.UserAdmin{
		Profile: &profile.Profile{
			User: &profile.User{},
		},
	}

	err = a.db.QueryRow(ctx, query, args...).
		Scan(
			&u.Profile.User.ID,
			&u.Profile.User.Username,
			&u.Profile.Email,
			&u.Profile.User.AvatarID,
			&u.Profile.User.Rating,
			&u.Profile.Coins,
			&u.Profile.User.CreatedAt,
			&u.Profile.User.LastLoginAt,
			&u.DeletedAt,
		)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("user", "email", email)
	case err != nil:
		return nil, apperrors.Internal(err)
	}

	return &u, nil
}

func (a *Admin) BanUser(ctx context.Context, userID string) error {
	builder := dbx.StatementBuilder.
		Update("users").
		Set("deleted_at", time.Now()).
		Where(squirrel.Eq{"id": userID})

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

func (a *Admin) UnbanUser(ctx context.Context, userID string) error {
	builder := dbx.StatementBuilder.
		Update("users").
		Set("deleted_at", nil).
		Where(squirrel.Eq{"id": userID})

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
