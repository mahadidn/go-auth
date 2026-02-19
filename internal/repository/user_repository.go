package repository

import (
	"context"
	"database/sql"
	"errors"
	"golang-auth/internal/domain"

	"github.com/google/uuid"
)

// definisikan struct secara private
type userRepository struct {
	db *sql.DB
}


// buat constructor
func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

// buat
func (u *userRepository) Create(ctx context.Context, user *domain.User) error {
	
	query := `INSERT INTO users (id, username, email, password, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?)`

	// ubah uuid ke format mysql
	idBytes, err := user.ID.MarshalBinary()
	if err != nil {
		return  err
	}

	// eksekusi query sql menggunakan ExecContext
	_, err = u.db.ExecContext(ctx, query,
		idBytes,
		user.Username,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)

	// return err, biar service tau kalau ini return error atau nil, biar lebih ringkas juga kodenya
	return  err
}

// cari berdasarkan id
func (u *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	
	// query := `SELECT id, username, email, created_at, updated_at`
	
	user := &domain.User{}

	return user, nil
}

// cari berdasarkan email
func (u *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	
	query := `SELECT id, password FROM users WHERE email = ?`

	var user domain.User
	var binID []byte

	// gunakan queryrowcontext
	err := u.db.QueryRowContext(ctx, query, email).Scan(
		&binID,
		&user.Password,
	)
	// cek apakah user ditemukan
	if err != nil {
		// not found error
		if err == sql.ErrNoRows {
			return  nil, errors.New("User not found")
		}
		// error tak terduga
		return nil, err
	}

	// kalau data ditemukan konversi dari bytes ke uuid.UUID
	user.ID, err = uuid.FromBytes(binID)
	if err != nil {
		return nil, err
	}
	user.Email = email

	return  &user, nil
}

// update
func (u *userRepository) Update(ctx context.Context, user *domain.User) error {
	
	query := `UPDATE users SET username = ?, email = ?, updated_at = ? WHERE id = ?`


	// ubah uuid ke format mysql
	idBytes, err := user.ID.MarshalBinary()
	if err != nil {
		return  err
	}

	res, err := u.db.ExecContext(ctx, query, 
		user.Username,
		user.Email,
		user.UpdatedAt,
		idBytes,
	)

	if err == nil {
		rows, _ := res.RowsAffected()
		if rows == 0 {
			return  errors.New("No user updated")
		}
	}

	// return err, agar service tau jika terjadi error
	return  err
}

// delete
func (u *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	
	query := `DELETE from users where id = ?`

	// konversi uuid ke format yg dimengerti mysql 
	idBytes, err := id.MarshalBinary()
	if err != nil {
		return  err
	}

	res, err := u.db.ExecContext(ctx, query, idBytes)

	if err == nil {
		rows, _ := res.RowsAffected()
		if rows == 0 {
			return errors.New("No user found to delete")
		}
	}

	return  err
}

// tambah role
func (u *userRepository) AddRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	

	return  nil
}

// hapus role
func (u *userRepository) RemoveRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	

	return  nil
}

// ubah password
func (u *userRepository) ChangePassword(ctx context.Context, id uuid.UUID ,newPassword string) error {

	query := `UPDATE users SET password = ? WHERE id = ?`

	// ubah uuid ke format mysql
	idBytes, err := id.MarshalBinary()
	if err != nil {
		return  err
	}

	res, err := u.db.ExecContext(ctx, query, idBytes)

	if err == nil {
		rows, _ := res.RowsAffected()
		if rows == 0 {
			return  errors.New("No user password changed")
		}
	}

	return  err
}