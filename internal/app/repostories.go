package app

import (
	"log/slog"

	"github.com/Sanchir01/go-shortener/internal/feature/url"
	"github.com/Sanchir01/go-shortener/internal/feature/user"
	"github.com/Sanchir01/go-shortener/pkg/db"
)

type Repositories struct {
	UserRepository *user.Repository
	UrlRepository  *url.Repository
}

func NewRepositories(databases *db.Database, l *slog.Logger) *Repositories {
	return &Repositories{
		UserRepository: user.NewRepository(databases.PrimaryDB, l),
		UrlRepository:  url.NewRepository(databases.PrimaryDB, l),
	}
}
