package app

import (
	"log/slog"

	"github.com/Sanchir01/go-shortener/internal/feature/url"
	"github.com/Sanchir01/go-shortener/internal/feature/user"
)

type Handlers struct {
	UserHandler *user.Handler
	UrlHandler  *url.Handler
}

func NewHandlers(services *Services, l *slog.Logger) *Handlers {
	return &Handlers{
		UserHandler: user.NewHandler(services.UserService, l),
		UrlHandler:  url.NewHandler(services.UrlService, l),
	}
}
