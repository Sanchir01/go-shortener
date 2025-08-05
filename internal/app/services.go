package app

import (
	"log/slog"

	"github.com/Sanchir01/go-shortener/internal/feature/url"
	"github.com/Sanchir01/go-shortener/internal/feature/user"
	"github.com/Sanchir01/go-shortener/pkg/db"
)

type Services struct {
	UserService *user.Service
	UrlService  *url.Service
}

func NewServices(repo *Repositories, db *db.Database, l *slog.Logger) *Services {
	return &Services{
		UserService: user.NewService(repo.UserRepository, db.PrimaryDB, l),
		UrlService:  url.NewService(repo.UrlRepository, db.PrimaryDB, l),
	}
}
