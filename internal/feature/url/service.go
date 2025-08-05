package url

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Sanchir01/go-shortener/internal/domain/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UrlService interface {
	CreateUrl(ctx context.Context, userId uuid.UUID, url, alise string, tx pgx.Tx) error
	GetUrlByUserId(ctx context.Context, userId uuid.UUID) ([]models.Url, error)
	GetAllUrl(ctx context.Context) ([]models.Url, error)
}
type Service struct {
	repo      UrlService
	primaryDB *pgxpool.Pool
	l         *slog.Logger
}

func NewService(repo UrlService, primaryDB *pgxpool.Pool, l *slog.Logger) *Service {
	return &Service{
		repo:      repo,
		primaryDB: primaryDB,
		l:         l,
	}
}

func (s *Service) GetAllUrl(ctx context.Context) ([]models.Url, error) {
	const op = "Url.Service.GetAllUrl"
	log := s.l.With(slog.String("op", op))

	urls, err := s.repo.GetAllUrl(ctx)
	if err != nil {
		log.Error("msg string", err.Error())
		return nil, err
	}
	log.Info("Getting all URLs completed service")
	return urls, nil
}

func (s *Service) CreateUrl(ctx context.Context, userId uuid.UUID, url string) error {
	const op = "Url.Service.CreateUrl"
	log := s.l.With(slog.String("op", op))

	conn, err := s.primaryDB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Error("tx error", err.Error())
		return err
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = errors.Join(err, rollbackErr)
				log.Error("rollback error", rollbackErr.Error())
				return
			}
		}
	}()

	if err := s.repo.CreateUrl(ctx, userId, url, "", tx); err != nil {
		log.Error("create url error", err.Error())
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("commit error", "msg", err.Error())
		return err
	}

	log.Info("Creating URL completed service")
	return nil
}
