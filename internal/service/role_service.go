package service

import (
	"context"
	"database/sql"
	"golang-auth/internal/domain"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type roleService struct {
	roleRepository domain.RoleRepository
	db *sql.DB
	validate *validator.Validate
}

func NewRoleService(roleRepository domain.RoleRepository, db *sql.DB, validate *validator.Validate) domain.RoleService {
	return &roleService{
		roleRepository: roleRepository,
		db: db,
		validate: validate,
	}
}


func (service *roleService) Create(ctx context.Context, reqRole domain.RoleCreateRequest) error {

	err := service.validate.Struct(reqRole)
	if err != nil {
		return err
	}

	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// aktifkan transaksi pada repository
	repoTx := service.roleRepository.WithTx(tx)

	uuid7, _ := uuid.NewV7()
	now := time.Now()

	err = repoTx.Create(ctx, &domain.Role{
		ID: uuid7,
		Name: reqRole.Name,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return err
	}

	// assign permission ke role
	err = repoTx.AssignPermission(ctx, uuid7, reqRole.PermissionIDs)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (service *roleService) Update(ctx context.Context, req domain.RoleUpdateRequest) error {
	err := service.validate.Struct(req)
	if err != nil {
		return err
	}

	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	repoTx := service.roleRepository.WithTx(tx)

	now := time.Now()

	// update data role
	err = repoTx.Update(ctx, &domain.Role{
		ID: req.ID,
		Name: req.Name,
		UpdatedAt: now,
	})
	if err != nil {
		return err
	}

	// update data permission
	// 1. remove
	err = repoTx.RemoveAllPermissions(ctx, req.ID)
	if err != nil {
		return err
	}
	// 2. attach
	err = repoTx.AssignPermission(ctx, req.ID, req.PermissionIDs)
	if err != nil {
		return err
	}
	
	return tx.Commit()
}

func (service *roleService) FindById(ctx context.Context, id uuid.UUID) (*domain.RoleWithUsersAndPermissions, error) {
	return nil, nil
}

func (service *roleService) FindAll(ctx context.Context) ([]domain.Role, error) {
	return nil, nil
}

func (service *roleService) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}