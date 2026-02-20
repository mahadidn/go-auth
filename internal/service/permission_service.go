package service

import (
	"context"
	"database/sql"
	"golang-auth/internal/domain"

	"github.com/google/uuid"
)

type permissionService struct {
	permissionRepository domain.PermissionRepository
	db *sql.DB
}

func NewPermissionService(permissionRepository domain.PermissionRepository, db *sql.DB) domain.PermissionService {
	return &permissionService{
		permissionRepository: permissionRepository,
		db: db,
	}
}

func (service permissionService) FindAll(ctx context.Context) ([]domain.Permission, error) {
	return service.permissionRepository.FindAll(ctx)
}

func (service permissionService) GetPermissionsByUserID(ctx context.Context, userID uuid.UUID) ([]string, error) {
	return service.permissionRepository.GetPermissionsByUserID(ctx, userID)
}

