package domain

import (
	"context"
	"database/sql"
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

// DTO untuk request create, update
type RoleCreateRequest struct {
	Name		  string		`validate:"required,min=3,max=100"`		
	PermissionIDs []uuid.UUID	`validate:"required,dive,uuid7"`
}

type RoleUpdateRequest struct {
	ID		uuid.UUID			`validate:"required,uuid7"`
	Name	string				`validate:"required,min=3,max=100"`
	PermissionIDs []uuid.UUID	`validate:"required,dive,uuid7"`
}

// repository interface
type RoleRepository interface {
	Create(ctx context.Context, r *Role) error
	AssignPermission(ctx context.Context, roleID uuid.UUID, permID []uuid.UUID) error
    RemoveAllPermissions(ctx context.Context, roleID uuid.UUID) error
	FindById(ctx context.Context, id uuid.UUID) (*RoleWithUsersAndPermissions, error)
	FindAll(ctx context.Context) ([]Role, error)
	Update(ctx context.Context, r *Role) error
	Delete(ctx context.Context, id uuid.UUID) error

	WithTx(tx *sql.Tx) RoleRepository
}

// service interface
type RoleService interface {
	Create(ctx context.Context, req RoleCreateRequest) error
	Update(ctx context.Context, req RoleUpdateRequest) error
	FindById(ctx context.Context, id uuid.UUID) (*RoleWithUsersAndPermissions, error)
	FindAll(ctx context.Context) ([]Role, error)
	Delete(ctx context.Context, id uuid.UUID) error
}