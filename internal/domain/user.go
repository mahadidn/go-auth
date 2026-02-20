package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID    `json:"id"`
	Username    string 	     `json:"username"`
	Email       string 	     `json:"email"`
	Password    string 	     `json:"-"`
	Roles       []Role       `json:"roles,omitempty"`
	Permissions []Permission `json:"permissions,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// DTO (Data Transfer Objects)
type UserCreateRequest struct {
	Username	string		`json:"username" validate:"required,min=3,max=50"`
	Email		string  	`json:"email" validate:"required,email"`
	Password	string  	`json:"password" validate:"required,min=6"`
	RoleIDs		[]uuid.UUID `json:"role_ids" validate:"omitempty,dive,uuid7"`
	PermissionIDs []uuid.UUID `json:"permission_ids" validate:"omitempty,dive,uuid7"`
}

type UserUpdateRequest struct {
	ID		 uuid.UUID   `json:"-"`
	Username string	     `json:"username" validate:"required,min=3,max=50"`
	Email 	 string		 `json:"email" validate:"required,email"`
	RoleIDs  []uuid.UUID `json:"role_ids" validate:"omitempty,dive,uuid7"`
	PermissionIDs []uuid.UUID `json:"permission_ids" validate:"omitempty,dive,uuid7"`
}

type UserChangePasswordRequest struct {
    OldPassword     string `json:"old_password" validate:"required"`
    NewPassword     string `json:"new_password" validate:"required,min=6,nefield=OldPassword"`
    
    // eqfield=NewPassword memastikan input ini sama persis ketikannya dengan NewPassword
    ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

type UserRepository interface {
	Create(ctx context.Context, u *User) error
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id uuid.UUID) error

	// role management
	AssignRoles(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) error
	RemoveAllRoles(ctx context.Context, userID uuid.UUID) error

	// permission management
	AssignPermissions(ctx context.Context, userID uuid.UUID, permissionIDs []uuid.UUID) error
	RemoveAllPermissions(ctx context.Context, userID uuid.UUID) error

	// password management
	ChangePassword(ctx context.Context, id uuid.UUID, newPassword string) error

	WithTx(tx *sql.Tx) UserRepository
}


type UserService interface {
	Create(ctx context.Context, req UserCreateRequest) error
	Update(ctx context.Context, req UserUpdateRequest) error
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Delete(ctx context.Context, id uuid.UUID) error

	ChangePassword(ctx context.Context, id uuid.UUID, req UserChangePasswordRequest) error	
}