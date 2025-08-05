package url

import (
	"github.com/Sanchir01/currency-wallet/pkg/api"
	"github.com/Sanchir01/go-shortener/internal/domain/models"
)

type GetAllUrlResponse struct {
	Response api.Response
	Urls     []models.Url `json:"urls"`
}
