package url

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Sanchir01/currency-wallet/pkg/api"
	"github.com/Sanchir01/go-shortener/internal/domain/models"
	"github.com/go-chi/render"
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
