package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID   	  uuid.UUID `json:"id"`
	Name 	  string	`json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RoleWithUsersAndPermissions struct {
	ID 			uuid.UUID 	 `json:"id"`
	Name 		string 		 `json:"name"`
	Users 		[]User 		 `json:"users"`
	Permissions []Permission `json:"permissions"`
}

type RoleRepository interface {
	Create(ctx context.Context, r *Role) error
	FindById(ctx context.Context, id uuid.UUID) (*RoleWithUsersAndPermissions, error)
	FindAll(ctx context.Context) ([]Role, error)
	Update(ctx context.Context, r *Role) error
	Delete(ctx context.Context, id uuid.UUID) error
}