package url

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Sanchir01/currency-wallet/pkg/api"
	contextkey "github.com/Sanchir01/go-shortener/internal/domain/constants"
	"github.com/Sanchir01/go-shortener/internal/domain/models"
	"github.com/Sanchir01/go-shortener/internal/feature/user"
	"github.com/Sanchir01/go-shortener/pkg/logger"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=UrlHandler
type UrlHandler interface {
	GetAllUrl(ctx context.Context) ([]models.Url, error)
	CreateUrl(ctx context.Context, userId uuid.UUID, url string) error
}
type Handler struct {
	service *Service
	l       *slog.Logger
}

func NewHandler(service *Service, l *slog.Logger) *Handler {
	return &Handler{
		service: service,
		l:       l,
	}
}

// @Summary  GetAllUrlHandler
// @Tags url
// @Description Get all urls admin
// @Accept json
// @Produce json
// @Success 200 {object}  GetAllUrlResponse
// @Failure 400,404 {object}  api.Response
// @Failure 500 {object}  api.Response
// @Router /url [get]
func (h *Handler) GetAllUrlHandler(w http.ResponseWriter, r *http.Request) {
	const op = "Url.Handler.GetAllUrl"
	log := h.l.With(slog.String("op", op))
	urls, err := h.service.GetAllUrl(r.Context())
	if err != nil {
		log.Error("error", "msg", err.Error())
		return
	}
	log.Info("getting all urls repo")
	fmt.Println("getting all urls repo", urls)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, GetAllUrlResponse{
		Response: api.OK(),
		Urls:     urls,
	})
}

// @Summary  CreateUrlHandler
// @Tags url
// @Description Create url
// @Accept json
// @Produce json
// @Param input body CreateUrlRequest true "login body"
// @Success 200 {object}  GetAllUrlResponse
// @Failure 400,404 {object}  api.Response
// @Failure 500 {object}  api.Response
// @Router /url/save [post]
func (h *Handler) CreateUrlHandler(w http.ResponseWriter, r *http.Request) {
	const op = "Url.Handler.CreateUrl"
	log := h.l.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	var req CreateUrlRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request body", slog.Any("err", err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("Ошибка при валидации данных"))
		return
	}
	if err := validator.New().Struct(req); err != nil {
		log.Error("invalid request", logger.Err(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid request"))
		return
	}
	claims, ok := r.Context().Value(contextkey.UserIDCtxKey).(*user.Claims)
	if !ok {
		log.Error("failed to parse product uuid")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, api.Error("failed send coins"))
		return
	}
	if err := h.service.CreateUrl(r.Context(), claims.ID, req.Url); err != nil {
		log.Error("failed to create url", logger.Err(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("failed to create url"))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, CreateUrlResponse{
		Response: api.OK(),
	})
}
