package models

import (
	"time"

	"github.com/google/uuid"
)

type Url struct {
	ID        uuid.UUID `db:"id"`
	Alias     string    `db:"alias"`
	Url       string    `db:"url"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	UserID    uuid.UUID `db:"user_id"`
}
