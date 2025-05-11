package dto

import (
	"github.com/google/uuid"
	"time"
)

type RegisterRequest struct {
	Id        uuid.UUID  `json:"id"`
	Email     string     `json:"email"`
	Hash      string     `json:"hash"`
	CreatedAt *time.Time `json:"created_at"`
}
