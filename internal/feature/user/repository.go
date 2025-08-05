package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/Sanchir01/currency-wallet/pkg/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	primaryDB *pgxpool.Pool
	l         *slog.Logger
}

func NewRepository(primaryDB *pgxpool.Pool, l *slog.Logger) *Repository {
	return &Repository{primaryDB: primaryDB, l: l}
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*DatabaseUser, error) {
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query, arg, err := sq.
		Select("id, email,title, version,password").
		From("public.users").
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	var userDB DatabaseUser
	if err := conn.QueryRow(ctx, query, arg...).Scan(&userDB.ID, &userDB.Email, &userDB.Name, &userDB.Version, &userDB.Password); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorUserNotFound
		}
		return nil, err
	}
	return &userDB, nil
}

func (r *Repository) CreateUser(ctx context.Context, email, username string, password []byte, tx pgx.Tx) (*uuid.UUID, error) {
	query, arg, err := sq.
		Insert("users").
		Columns("email", "password", "title").
		Values(email, password, username).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, utils.ErrorQueryString
	}
	var id uuid.UUID

	if err := tx.QueryRow(ctx, query, arg...).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, fmt.Errorf("пользователь с таким email или username уже существует")
			}
		}

		return nil, err
	}
	return &id, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*DatabaseUser, error) {
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query, arg, err := sq.
		Select("id, email, title, version").
		From("public.users").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	var userDB DatabaseUser
	if err := conn.QueryRow(ctx, query, arg...).Scan(&userDB.ID, &userDB.Email, &userDB.Version); err != nil {
		return nil, err
	}
	return &userDB, nil
}
