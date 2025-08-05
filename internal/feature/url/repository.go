package url

import (
	"context"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/Sanchir01/go-shortener/internal/domain/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	primaryDB *pgxpool.Pool
	l         *slog.Logger
}

func NewRepository(primaryDB *pgxpool.Pool, l *slog.Logger) *Repository {
	return &Repository{
		primaryDB: primaryDB,
		l:         l,
	}
}

func (r *Repository) CreateUrl(ctx context.Context, userId uuid.UUID, url, alise string, tx pgx.Tx) error {
	const op = "Url.Repository.CreateUrl"
	log := r.l.With(slog.String("op", op))

	query, args, err := sq.
		Insert("url").
		Columns("user_id", "url", "alias").
		Values(userId, url, alise).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		log.Error("error", err.Error())
		return err
	}

	id := uuid.New()
	args = append(args, id.String())

	if _, err := tx.Exec(ctx, query, args...); err != nil {
		log.Error("error", err.Error())
		return err
	}

	log.Info("creating url repo")
	return nil
}

func (r *Repository) GetUrlByUserId(ctx context.Context, userId uuid.UUID) ([]models.Url, error) {
	const op = "Url.Repository.GetUrlByUserId"
	log := r.l.With(slog.String("op", op))
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		log.Error("error", err.Error())
		return nil, err
	}
	defer conn.Release()

	query, args, err := sq.
		Select("id, url,alias, created_at, updated_at").
		From("url").
		Where(sq.Eq{"user_id": userId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		log.Error("error", err.Error())
		return nil, err
	}

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		log.Error("error", err.Error())
		return nil, err
	}

	urls := make([]models.Url, 0, 100)
	defer rows.Close()
	for rows.Next() {
		var url models.Url
		if err := rows.Scan(&url.ID, &url.Url, &url.Alias, &url.CreatedAt, &url.UpdatedAt); err != nil {
			log.Error("error", err.Error())
			return nil, err
		}
		urls = append(urls, url)
	}
	if err := rows.Err(); err != nil {
		log.Error("error", err.Error())
		return nil, err
	}
	log.Info("getting all urls repo")
	return urls, nil
}

func (r *Repository) GetAllUrl(ctx context.Context) ([]models.Url, error) {
	const op = "Url.Repository.GetAllUrl"
	log := r.l.With(slog.String("op", op))
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		log.Error("error", err.Error())
		return nil, err
	}
	defer conn.Release()

	query, args, err := sq.
		Select("id, url,alias, created_at, updated_at").
		From("url").
		ToSql()

	if err != nil {
		log.Error("error", err.Error())
		return nil, err
	}

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		log.Error("error", err.Error())
		return nil, err
	}

	urls := make([]models.Url, 0, 100)
	defer rows.Close()
	for rows.Next() {
		var url models.Url
		if err := rows.Scan(&url.ID, &url.Url, &url.Alias, &url.CreatedAt, &url.UpdatedAt); err != nil {
			log.Error("error", err.Error())
			return nil, err
		}
		urls = append(urls, url)
	}
	if err := rows.Err(); err != nil {
		log.Error("error", err.Error())
		return nil, err
	}
	log.Info("getting all urls repo")
	return urls, nil
}
