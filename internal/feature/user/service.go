package user

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Sanchir01/currency-wallet/pkg/utils"
	"github.com/Sanchir01/go-shortener/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	repository ServiceUser
	log        *slog.Logger
	primaryDB  *pgxpool.Pool
}

type ServiceUser interface {
	CreateUser(ctx context.Context, email, username string, password []byte, tx pgx.Tx) (*uuid.UUID, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*DatabaseUser, error)
	GetUserByEmail(ctx context.Context, email string) (*DatabaseUser, error)
}

func NewService(r ServiceUser, db *pgxpool.Pool, l *slog.Logger) *Service {
	return &Service{
		repository: r,
		primaryDB:  db,
		log:        l,
	}
}

func (s *Service) Register(ctx context.Context, email, username, password string) (*uuid.UUID, error) {
	const op = "User.Service.Register"
	log := s.log.With(slog.String("op", op))
	conn, err := s.primaryDB.Acquire(ctx)
	if err != nil {

		return nil, err
	}
	defer conn.Release()
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Error("tx error", logger.Err(err))
		return nil, err
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
	hashedPassword, err := GeneratePasswordHash(password)
	if err != nil {
		log.Error("error generating password hash", logger.Err(err))
		return nil, err
	}
	user, err := s.repository.CreateUser(ctx, email, username, hashedPassword, tx)
	if err != nil {
		log.Error("error creating user", logger.Err(err))
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("tx commit error", logger.Err(err))
	}
	log.Info("user created success", logger.Err(err))
	return user, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*DatabaseUser, error) {
	const op = "User.Service.Login"
	log := s.log.With(slog.String("op", op))
	user, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error("error getting user by email", logger.Err(err))
		return nil, err
	}
	ok := VerifyPassword(user.Password, password)
	if !ok {
		log.Error("invalid password")
		return nil, utils.ErrorInvalidPassword
	}
	log.Info("user service logged in user")
	return user, nil
}
