package url

import (
	"github.com/Sanchir01/currency-wallet/pkg/api"
	"github.com/Sanchir01/go-shortener/internal/domain/models"
)

type GetAllUrlResponse struct {
	api.Response
	Urls []models.Url `json:"urls"`
}
type CreateUrlResponse struct {
	api.Response
	Url string `json:"url"`
}
type CreateUrlRequest struct {
	Url string `json:"url" validate:"required"`
}
