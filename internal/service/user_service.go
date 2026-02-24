package service

import (
	"context"
	"database/sql"
	"errors"
	"golang-auth/internal/domain"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepository domain.UserRepository
	db *sql.DB
	validate *validator.Validate
}

func NewUserService(userRepository domain.UserRepository, db *sql.DB, validate *validator.Validate) domain.UserService {
	return &userService{
		userRepository: userRepository,
		db: db,
		validate: validate,
	}
}

func (service *userService) Create(ctx context.Context, req domain.UserCreateRequest) error {
	
	err := service.validate.Struct(req)
	if err != nil {
		return err
	}

	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	repoTx := service.userRepository.WithTx(tx)

	// cek email sudah dipakai atau belum
	_, err = repoTx.FindByEmail(ctx, req.Email)
	if err == nil {
		return errors.New("Email sudah digunakan")
	}else if err != sql.ErrNoRows {
		// kalau errornya bukan karna data kosong (misal DB mati), kembalikan errornya
		return err
	}

	uuid7, _ := uuid.NewV7()
	now := time.Now()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err // gagal melakukan hashing
	}

	// buat user
	err = repoTx.Create(ctx, &domain.User{
		ID: uuid7,
		Username: req.Username,
		Email: req.Email,
		Password: string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return err
	}

	// assign role (kalau role-nya ada isinya jalanin querynya)
	if len(req.RoleIDs) > 0 {
		err = repoTx.AssignRoles(ctx, uuid7, req.RoleIDs)
		if err != nil {
			return err
		}
	}

	// assign permission (kalau permission-nya ada isinya jalanin querynya)
	if len(req.PermissionIDs) > 0 {
		err = repoTx.AssignPermissions(ctx, uuid7, req.PermissionIDs)
		if err != nil {
			return err
		}
	}
	
	return tx.Commit()
}

func (service *userService) Update(ctx context.Context, req domain.UserUpdateRequest) error {
	err := service.validate.Struct(req)
	if err != nil {
		return err
	}
	
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	repoTx := service.userRepository.WithTx(tx)

	// cek email, email yg dirubah sudah dipakai orang belum
	existingUser, err := repoTx.FindByEmail(ctx, req.Email)
	if err == nil {
		// email ketemu, cek apakah ID-nya beda dengan ID user yg sedang di-update
		if existingUser.ID != req.ID {
			return errors.New("email sudah digunakan oleh pengguna lain")
		}
	}else if err != sql.ErrNoRows {
		// jika errornya bukan karna data kosong, misal koneksi db putus maka batalkan
		return err
	}

	now := time.Now()
	err = repoTx.Update(ctx, &domain.User{
		ID: req.ID,
		Username: req.Username,
		Email: req.Email,
		UpdatedAt: now,
	})

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (service *userService) AssignRoles(ctx context.Context, id uuid.UUID, req domain.AssignRoleRequest) error {
	
	err := service.validate.Struct(req)
	if err != nil {
		return err
	}

	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	repoTx := service.userRepository.WithTx(tx)
	
	// remove all roles
	err = repoTx.RemoveAllRoles(ctx, id)
	if err != nil {
		return err
	}

	// add all roles
	err = repoTx.AssignRoles(ctx, id, req.RoleIDs)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (service *userService) AssignPermissions(ctx context.Context, id uuid.UUID, req domain.AssignPermissionRequest) error {
	err := service.validate.Struct(req)
	if err != nil {
		return err
	}

	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	repoTx := service.userRepository.WithTx(tx)
	
	// remove all permission
	err = repoTx.RemoveAllPermissions(ctx, id)
	if err != nil {
		return err
	}
	
	// add all permission
	err = repoTx.AssignPermissions(ctx, id, req.PermissionIDs)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (service *userService) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	res, err := service.userRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (service *userService) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	res, err := service.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (service *userService) Delete(ctx context.Context, id uuid.UUID) error {
	err := service.userRepository.Delete(ctx, id)
	return err
}

func (service *userService) ChangePassword(ctx context.Context, id uuid.UUID, req domain.UserChangePasswordRequest) error {
	
	err := service.validate.Struct(req)
	if err != nil {
		return err
	}

	user, err := service.userRepository.FindByID(ctx, id)
	if err != nil {
		return err // error user tidak ditemukan
	}

	// cocokkan OldPassword (input user) dengan password di DB
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		return errors.New("password lama tidak sesuai")
	}

	// hash new password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = service.userRepository.ChangePassword(ctx, id, string(hashedPassword))
	
	return err
}
