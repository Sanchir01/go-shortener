package user

import (
	"time"

	"github.com/Sanchir01/go-shortener/pkg/api"
	"github.com/google/uuid"
)

type DatabaseUser struct {
	ID        uuid.UUID `db:"id"`
	Email     string    `db:"email"`
	Name      string    `db:"username"`
	Password  []byte    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Role      string    `db:"role"`
	Version   int64     `db:"version"`
}
type AuthRequest struct {
	Email    string `json:"email" validate:"required"`
	Username string `json:"title" validate:"required,min=1,max=100"`
	Password string `json:"password" validate:"required,min=6"`
}

type GoogleRegisterRequest struct {
	Code string `json:"code" validate:"required"`
}
type GoogleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token"`
}
type AuthResponse struct {
	api.Response
	Email    string
	Username string
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	api.Response
	Email    string `json:"email"`
	Username string `json:"username" `
}
