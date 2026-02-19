package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID   	  uuid.UUID `json:"id"`
	Name 	  string	`json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RoleWithUsers struct {
    ID    uuid.UUID `json:"id"`
    Name  string    `json:"name"`
    Users []User    `json:"users"` 
}

type RoleWithPermissions struct {
	ID 			uuid.UUID	 `json:"id"`
	Name 		string		 `json:"name"`
	Permissions []Permission `json:"permissions"`
}