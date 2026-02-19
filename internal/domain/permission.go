package domain

import (
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID 		  uuid.UUID `json:"id"`
	Name 	  string 	`json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PermissionWithUsers struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Users []User    `json:"users"`
}