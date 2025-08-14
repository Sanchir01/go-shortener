package app

import (
	"context"
	"fmt"
	"log/slog"

	tgbot "github.com/Sanchir01/go-shortener/internal/bot"
	"github.com/Sanchir01/go-shortener/internal/config"
	httpserver "github.com/Sanchir01/go-shortener/internal/server/http"
	"github.com/Sanchir01/go-shortener/pkg/db"
	"github.com/Sanchir01/go-shortener/pkg/logger"
)

type App struct {
	Log              *slog.Logger
	Cfg              *config.Config
	Handlers         *Handlers
	Services         *Services
	Bot              *tgbot.TGBot
	DB               *db.Database
	HttpServer       *httpserver.Server
	PrometheusServer *httpserver.Server
}

func NewApp(ctx context.Context) (*App, error) {
	cfg := config.InitConfig()
	fmt.Println(cfg)
	l, _ := logger.SetupLogger(ctx, cfg.Env)
	database, err := db.NewDataBases(cfg, ctx, l)
	if err != nil {
		l.Error("db error", slog.String("error", err.Error()))
		return nil, err
	}
	repo := NewRepositories(database, l)
	services := NewServices(repo, database, l)

	botInstance, err := tgbot.New(ctx, services.UrlService, services.UserService, l)
	if err != nil {
		l.Error("bot connect error:", slog.String("error", err.Error()))
		return nil, err
	}

	handlers := NewHandlers(services, l)
	httpServer := httpserver.NewHTTPServer(cfg.HttpServer.Host, cfg.HttpServer.Port, cfg.HttpServer.Timeout, cfg.HttpServer.IdleTimeout)
	prometheusServer := httpserver.NewHTTPServer(cfg.Prometheus.Host, cfg.Prometheus.Port, cfg.Prometheus.Timeout, cfg.Prometheus.IdleTimeout)

	return &App{
		Log:              l,
		Cfg:              cfg,
		HttpServer:       httpServer,
		Services:         services,
		Bot:              botInstance,
		Handlers:         handlers,
		DB:               database,
		PrometheusServer: prometheusServer,
	}, nil
}
