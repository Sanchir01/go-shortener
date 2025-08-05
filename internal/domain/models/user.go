package models

import (
	"time"

	contextkey "github.com/Sanchir01/go-shortener/internal/domain/constants"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID           `json:"id"`
	Name      string              `json:"username"`
	Email     string              `json:"email"`
	Password  string              `json:"password"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	Version   uint64              `json:"version"`
	Role      contextkey.UserRole `json:"role"`
}
