package repository

import (
	"context"
	"database/sql"
	"errors"
	"golang-auth/internal/domain"
	"strings"

	"github.com/google/uuid"
)

type roleRepository struct {
	db DBTX
}

func NewRoleRepository(db *sql.DB) domain.RoleRepository {
	return &roleRepository{
		db: db,
	}
}

func (repo *roleRepository) WithTx(tx *sql.Tx) domain.RoleRepository {
	return &roleRepository{
		db: tx,
	}
}

func (repo *roleRepository) Create(ctx context.Context, role *domain.Role) error {
	
	query := `INSERT INTO roles (id, name, created_at, updated_at) 
			  VALUES (?, ?, ?, ?)`

	// ubah uuid ke format mysql
	idBytes, err := role.ID.MarshalBinary()
	if err != nil {
		return err
	}

	// eksekusi query SQL menggunakan execcontext
	_, err = repo.db.ExecContext(ctx, query,
		idBytes,
		role.Name,
		role.CreatedAt,
		role.UpdatedAt,
	)

	return err
}

// buat get role berdasarkan ID
func (repo *roleRepository) FindById(ctx context.Context, id uuid.UUID) (*domain.RoleWithUsersAndPermissions, error) {
	
	res := &domain.RoleWithUsersAndPermissions{}

	// konversi uuid ke bytes
	binID, _ := id.MarshalBinary()

	// query pertama: ambil data role
	queryRole := `SELECT id, name FROM roles WHERE id = ?`
	var roleBinID []byte
	err := repo.db.QueryRowContext(ctx, queryRole, binID).Scan(
		&roleBinID,
		&res.Name,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	res.ID, _ = uuid.FromBytes(roleBinID)

	// query kedua: ambil permissions
	queryPermission := `SELECT p.id, p.name FROM permissions as p
						 JOIN role_has_permissions as rhp ON p.id = rhp.permission_id
						 WHERE rhp.role_id = ?`
	
	rowsP, err := repo.db.QueryContext(ctx, queryPermission, binID)
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
		// konversi id ke uuid
		permission.ID, _ = uuid.FromBytes(permissionBinId)

		// append permission ke res
		res.Permissions = append(res.Permissions, permission)
	}
	if err = rowsP.Err(); err != nil {
		return nil, err
	}

	// query ketiga: Ambil users
	queryUser := `SELECT u.id, u.username, u.email FROM users as u
				  JOIN user_has_roles as uhr ON u.id = uhr.user_id
				  WHERE uhr.role_id = ?`
	
	rowsU, err := repo.db.QueryContext(ctx, queryUser, binID)
	if err != nil {
		return nil, err
	}
	defer rowsU.Close()

	for rowsU.Next() {
		var user domain.User
		var userBinId []byte
		err := rowsU.Scan(
			&userBinId,
			&user.Username,
			&user.Email,
		)
		if err != nil {
			return nil, err
		}
		// konversi id ke uuid
		user.ID, _ = uuid.FromBytes(userBinId)

		// append users ke res
		res.Users = append(res.Users, user)
	}
	if err = rowsU.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (repo *roleRepository) FindAll(ctx context.Context) ([]domain.Role, error) {
	
	query := `SELECT id, name, created_at, updated_at FROM roles`
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var roles []domain.Role
	for rows.Next() {
		var role domain.Role
		var binID []byte

		err := rows.Scan(
			&binID,
			&role.Name,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		// konversi dari bytes ke uuid.UUID
		role.ID, err = uuid.FromBytes(binID)
		if err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

func (repo *roleRepository) GetRoleByUserID(ctx context.Context, userID uuid.UUID) ([]string, error) {
	
	userBinID, _ := userID.MarshalBinary()
	
	queryRole := `SELECT r.name FROM roles as r
				  JOIN user_has_roles as uhr ON r.id = uhr.role_id
				  WHERE uhr.user_id = ?`
	
	rows, err := repo.db.QueryContext(ctx, queryRole, userBinID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		err := rows.Scan(
			&role,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (repo *roleRepository) Update(ctx context.Context, role *domain.Role) error {
	
	query := `UPDATE roles SET name = ?, updated_at = ? WHERE id = ?`

	// ubah uuid ke format mysql
	idBytes, err := role.ID.MarshalBinary()
	if err != nil {
		return err
	}

	res, err := repo.db.ExecContext(ctx, query,
		role.Name,
		role.UpdatedAt,
		idBytes,
	)
	if err == nil {
		rows, _ := res.RowsAffected()
		if rows == 0 {
			return errors.New("no role updated")
		}
	}
	
	return err
}

func (repo *roleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	
	query := `DELETE from roles WHERE id = ?`

	idBytes, err := id.MarshalBinary()
	if err != nil {
		return err
	}

	res, err := repo.db.ExecContext(ctx, query, idBytes)

	if err == nil {
		rows, _ := res.RowsAffected()
		if rows == 0 {
			return errors.New("no role found to delete") 
		}
	}
	
	return err
}

func (repo *roleRepository) AssignPermission(ctx context.Context, roleID uuid.UUID, permIDs []uuid.UUID) error {
	if len(permIDs) == 0 {
		return nil
	}

	query := `INSERT INTO role_has_permissions (role_id, permission_id) VALUES `

	// pre-allocation
	values := make([]interface{}, 0, len(permIDs) * 2) // *2 karena 1 role, 1 perm
	placeHolders := make([]string, 0, len(permIDs))

	roleBin, _ := roleID.MarshalBinary()

	// bangun string (?, ?) sebanyak jumlah permIDs
	for _, pID := range permIDs {
		placeHolders = append(placeHolders, "(?, ?)")

		permBin, _ := pID.MarshalBinary()
		values = append(values, roleBin, permBin)
	}

	// gabungkan query dasar dengan semua placeholder
	query += strings.Join(placeHolders, ", ")

	// eksekusi query
	_, err := repo.db.ExecContext(ctx, query, values...)

	return err
}

func (repo *roleRepository) RemoveAllPermissions(ctx context.Context, roleID uuid.UUID) error {
	query := `DELETE FROM role_has_permissions WHERE role_id = ?`

	roleBinID, _ := roleID.MarshalBinary()

	_, err := repo.db.ExecContext(ctx, query, roleBinID)

	return err
}