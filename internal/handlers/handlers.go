package httphandlers

import (
	"log/slog"
	"net/http"

	_ "github.com/Sanchir01/go-shortener/docs"
	"github.com/Sanchir01/go-shortener/internal/app"
	"github.com/Sanchir01/go-shortener/internal/handlers/customiddleware"
	"github.com/Sanchir01/go-shortener/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

func StartHTTTPHandlers(handlers *app.Handlers, domain string, l *slog.Logger) http.Handler {
	router := chi.NewRouter()
	custommiddleware(router, l)
	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", handlers.UserHandler.RegisterHandler)
			r.Post("/login", handlers.UserHandler.LoginHandler)
			r.Group(func(r chi.Router) {
				r.Use(customiddleware.AuthMiddleware(domain))
			})
			r.Route("/google", func(r chi.Router) {
				r.Get("/", handlers.UserHandler.GoogleLogin)
				r.Post("/register", handlers.UserHandler.GoogleCallback)
			})
		})
		r.Route("/url", func(r chi.Router) {
			r.Use(customiddleware.AuthMiddleware(domain))
			r.Post("/save", handlers.UrlHandler.CreateUrlHandler)
			r.Get("/", handlers.UrlHandler.GetAllUrlHandler)

		})
		r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, World!"))
		})
	})
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	return router
}
func custommiddleware(router *chi.Mux, l *slog.Logger) {
	router.Use(middleware.RequestID, middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(logger.NewMiddlewareLogger(l))
	router.Use(customiddleware.PrometheusMiddleware)
}
func StartPrometheusHandlers() http.Handler {
	router := chi.NewRouter()
	router.Handle("/metrics", promhttp.Handler())
	return router
}
