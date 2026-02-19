package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID 		  uuid.UUID `json:"id"`
	Name 	  string 	`json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PermissionRepository interface {
	FindAll(ctx context.Context) ([]Permission, error)
	// Ini yang akan dipakai oleh Middleware Routing nanti

    // Mengambil semua nama permission (misal: "user.create", "user.delete") 
    // baik dari Role maupun Direct Permission
	// return list permission aja
    GetPermissionsByUserID(ctx context.Context, userID uuid.UUID) ([]string, error)
	GetPermissionsByRoleID(ctx context.Context, roleID uuid.UUID) ([]string, error)
}