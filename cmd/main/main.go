package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sanchir01/go-shortener/internal/app"
	httphandlers "github.com/Sanchir01/go-shortener/internal/handlers"
)

// @title 🚀 URL-SHORTENER
// @version         1.0
// @description This is a sample server seller
// @termsOfService  http://swagger.io/terms/

// @host localhost:4200
// @BasePath /api/v1

// @securityDefinitions.apikey AccessTokenCookie
// @in cookie
// @name accessToken

// @securityDefinitions.apikey RefreshTokenCookie
// @in cookie
// @name refreshToken

// @contact.name GitHub
// @contact.url https://github.com/Sanchir01
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer cancel()
	application, err := app.NewApp(ctx)
	if err != nil {
		slog.Error("panic error", slog.Any("err", err))
		panic(err)
	}

	go func() {
		if err := application.HttpServer.Run(httphandlers.StartHTTTPHandlers(application.Handlers, application.Cfg.Domain, application.Log)); err != nil {
			application.Log.Error("Error while running http server", slog.String("error", err.Error()))
			cancel()
		}
	}()
	go func() {
		if err := application.PrometheusServer.Run(httphandlers.StartPrometheusHandlers()); err != nil {
			application.Log.Error("Error while running prometheus server", slog.String("error", err.Error()))
			cancel()
		}
	}()
	go func() {
		application.Bot.Bot.Start(ctx)
	}()
	<-ctx.Done()

	if err := application.HttpServer.Gracefull(ctx); err != nil {
		application.Log.Error("Close database", slog.String("error", err.Error()))
	}
	if err := application.DB.Close(); err != nil {
		application.Log.Error("Close database", slog.String("error", err.Error()))
	}
}
