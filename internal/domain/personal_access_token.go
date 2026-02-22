package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PersonalAccessToken struct {
    ID         uuid.UUID  `json:"id"`
	TokenHash  string 	  `json:"token"`
	UserID 	   uuid.UUID  `json:"user_id"`
	TokenName  string 	  `json:"token_name"`
	LastUsedAt *time.Time `json:"last_used_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

type PersonalAccessTokenRequest struct {
	UserID 	   uuid.UUID  `json:"user_id" validate:"required,uuid7"`
	TokenName      string 	  `json:"token" validate:"required,min=3,max=100"`
}


type PersonalAccessTokenRepository interface {
    // Digunakan saat Login
    Create(ctx context.Context, token *PersonalAccessToken) error
    
    // Digunakan oleh Middleware untuk cek apakah token valid/ada di DB
    FindByToken(ctx context.Context, token string) (*PersonalAccessToken, error)
    
    // Digunakan saat Logout (Hapus token ini saja)
    Delete(ctx context.Context, token string) error
    
    // Digunakan untuk "Logout from all devices"
    DeleteByUserID(ctx context.Context, userID uuid.UUID) error
    
    // Digunakan untuk update kolom last_used_at tiap kali user akses API
    UpdateLastUsed(ctx context.Context, token string) error
}

type PersonalAccessTokenService interface {
    Create(ctx context.Context, req PersonalAccessTokenRequest) (string, error)
    FindByToken(ctx context.Context, token string) (*PersonalAccessToken, error)
    Delete(ctx context.Context, token string) error
    DeleteByUserID(ctx context.Context, userID uuid.UUID) error
    UpdateLastUsed(ctx context.Context, token string) error
}