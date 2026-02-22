package repository

import (
	"context"
	"database/sql"
	"errors"
	"golang-auth/internal/domain"
	"strings"
	"time"

	"github.com/google/uuid"
)

// definisikan struct secara private
type userRepository struct {
	db DBTX
}


// buat constructor
func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) WithTx(tx *sql.Tx) domain.UserRepository {
	return &userRepository{
		db: tx,
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
	
	res := &domain.User{}
	
	// konversi uuid ke bytes
	binID, _ := id.MarshalBinary()
	
	// query pertama: ambil data user
	queryUser := `SELECT id, username, email, created_at, updated_at FROM users WHERE id = ?`
	var userBinId []byte
	err := u.db.QueryRowContext(ctx, queryUser, binID).Scan(
		&userBinId,
		&res.Username,
		&res.Email,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return  nil, err
	}
	res.ID, _ = uuid.FromBytes(userBinId)

	// query kedua: ambil role user
	queryRole := `SELECT r.id, r.name FROM roles as r
				  JOIN user_has_roles as uhr ON r.id = uhr.role_id
				  WHERE uhr.user_id = ?`
	rowsR, err := u.db.QueryContext(ctx, queryRole, binID)
	if err != nil {
		return nil, err
	}
	defer rowsR.Close()

	for rowsR.Next() {
		var role domain.Role
		var roleBinId []byte
		err := rowsR.Scan(
			&roleBinId,
			&role.Name,
		)
		if err != nil {
			return nil, err
		}
		// konversi id ke uuid
		role.ID, _ = uuid.FromBytes(roleBinId)

		// append roles ke user
		res.Roles = append(res.Roles, role)
	}
	if err = rowsR.Err(); err != nil {
		return nil, err
	}

	// query ketiga: ambil permission user
	queryPermission := `SELECT p.id, p.name FROM permissions as p
						JOIN user_has_permissions as uhp ON p.id = uhp.permission_id
						WHERE uhp.user_id = ?`
	rowsP, err := u.db.QueryContext(ctx, queryPermission, binID)
	if err != nil {
		return nil, err
	}
	defer rowsP.Close()

	for rowsP.Next() {
		var permission domain.Permission
		var permissionBinId []byte
		err := rowsP.Scan(
			&permissionBinId,
			&permission.Name,
		)
		if err != nil {
			return nil, err
		}
		permission.ID, _ = uuid.FromBytes(permissionBinId)

		res.Permissions = append(res.Permissions, permission)
	}
	if err = rowsP.Err(); err != nil {
		return nil, err
	}
	
	return res, nil
}

// cari berdasarkan email
func (u *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	
	query := `SELECT id, username, email, password FROM users WHERE email = ?`

	var user domain.User
	var binID []byte

	// gunakan queryrowcontext
	err := u.db.QueryRowContext(ctx, query, email).Scan(
		&binID,
		&user.Username,
		&user.Email,
		&user.Password,
	)
	// cek apakah user ditemukan
	if err != nil {
		// not found error
		if err == sql.ErrNoRows {
			return  nil, sql.ErrNoRows // kembalikan error aslinya agar mudah di cek
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
func (u *userRepository) AssignRoles(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) error {
	
	if len(roleIDs) == 0 {
		return nil
	}

	query := `INSERT INTO user_has_roles (user_id, role_id) VALUES `

	// konversi kedua ID
	userBinId, _ := userID.MarshalBinary()

	// pre-allocation slice
	values := make([]interface{}, 0, len(roleIDs)*2) // * 2 karna butuh user_id dan role_id
	placeHolders := make([]string, 0, len(roleIDs))

	for _, rID := range roleIDs {
		placeHolders = append(placeHolders, "(?, ?)")

		roleBinID, _ := rID.MarshalBinary()
		values = append(values, userBinId, roleBinID)
	}

	query += strings.Join(placeHolders, ", ")

	// eksekusi
	_, err := u.db.ExecContext(ctx, query, values...)

	return  err
}

// hapus role
func (u *userRepository) RemoveAllRoles(ctx context.Context, userID uuid.UUID) error {
	
	query := `DELETE FROM user_has_roles WHERE user_id = ?`

	userBinId, _ := userID.MarshalBinary()

	_, err := u.db.ExecContext(ctx, query, userBinId)

	return  err
}

// tambah permissions
func (u *userRepository) AssignPermissions(ctx context.Context, userID uuid.UUID, permissionIDs []uuid.UUID) error {
	
	if len(permissionIDs) == 0 {
		return nil
	}

	query := `INSERT INTO user_has_permissions (user_id, permission_id) VALUES `

	// konversi kedua ID
	userBinId, _ := userID.MarshalBinary()

	// pre-allocation slice
	values := make([]interface{}, 0, len(permissionIDs)*2) // * 2 karna butuh user_id dan permission_id
	placeHolders := make([]string, 0, len(permissionIDs))

	for _, pID := range permissionIDs {
		placeHolders = append(placeHolders, "(?, ?)")

		permBinID, _ := pID.MarshalBinary()
		values = append(values, userBinId, permBinID)
	}

	query += strings.Join(placeHolders, ", ")

	// eksekusi
	_, err := u.db.ExecContext(ctx, query, values...)

	return  err
}

// hapus permissions
func (u *userRepository) RemoveAllPermissions(ctx context.Context, userID uuid.UUID) error {
	
	query := `DELETE FROM user_has_permissions WHERE user_id = ?`

	userBinId, _ := userID.MarshalBinary()

	_, err := u.db.ExecContext(ctx, query, userBinId)

	return  err
}

// ubah password
func (u *userRepository) ChangePassword(ctx context.Context, id uuid.UUID ,newPassword string) error {

	query := `UPDATE users SET password = ?, updated_at = ? WHERE id = ?`

	// ubah uuid ke format mysql
	idBytes, err := id.MarshalBinary()
	if err != nil {
		return  err
	}

	res, err := u.db.ExecContext(ctx, query, 
		newPassword,
		time.Now(),
		idBytes,
	)

	if err == nil {
		rows, _ := res.RowsAffected()
		if rows == 0 {
			return  errors.New("No user password changed")
		}
	}

	return  err
}