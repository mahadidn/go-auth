package domain

import (
	"time"

	"github.com/google/uuid"
)

type PersonalAccessToken struct {
	Token      string 	  `json:"token"`
	UserID 	   uuid.UUID  `json:"user_id"`
	TokenName  string 	  `json:"token_name"`
	LastUsedAt *time.Time `json:"last_used_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
}
