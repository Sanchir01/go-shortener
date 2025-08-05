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
	app, err := app.NewApp(ctx)
	if err != nil {
		app.Log.Error("panic error", err.Error())
		panic(err)
	}

	go func() {
		if err := app.HttpServer.Run(httphandlers.StartHTTTPHandlers(app.Handlers, app.Cfg.Domain, app.Log)); err != nil {
			app.Log.Error("Error while running http server", slog.String("error", err.Error()))
			cancel()
		}
	}()
	go func() {
		if err := app.PrometheusServer.Run(httphandlers.StartPrometheusHandlers()); err != nil {
			app.Log.Error("Error while running prometheus server", slog.String("error", err.Error()))
			cancel()
		}
	}()
	<-ctx.Done()

	if err := app.HttpServer.Gracefull(ctx); err != nil {
		app.Log.Error("Close database", slog.String("error", err.Error()))
	}
	if err := app.DB.Close(); err != nil {
		app.Log.Error("Close database", slog.String("error", err.Error()))
	}
}
