package repository

import (
	"context"
	"database/sql"
	"fmt"
	"golang-auth/internal/domain"
	"strings"

	"github.com/google/uuid"
)

// definisikan struct secara private
type permissionRepository struct {
	db *sql.DB
}

// buat constructor
func NewPermissionRepository(db *sql.DB) domain.PermissionRepository {
	return &permissionRepository{
		db: db,
	}
}

// get all permissions
func(p *permissionRepository) FindAll(ctx context.Context) ([]domain.Permission, error)   {
	
	query := `SELECT id, name, created_at, updated_at FROM permissions`
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []domain.Permission
	for rows.Next() {
		var permission domain.Permission
		var binID []byte

		err := rows.Scan(
			&binID, 
			&permission.Name, 
			&permission.CreatedAt, 
			&permission.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// konversi dari bytes ke uuid.UUID
		permission.ID, err = uuid.FromBytes(binID)
		if err != nil {
			return  nil, err
		}

		permissions = append(permissions, permission)
	}

	// Cek apakah ada error selama proses looping
    if err = rows.Err(); err != nil {
        return nil, err
    }

	return permissions, nil
}

// get permission by user id
func(p *permissionRepository) GetPermissionsByUserID(ctx context.Context, userID uuid.UUID) ([]string, error) {
	// query pertama ke role dan ke user_has_roles, buat cek role si user ada akses ke permission-nya atau tidak (indirect)
	// query kedua ke user_has_permissions, buat cek si user punya akses langsung ke permissions atau tidak (direct) 
	// union buat ambil hasil query select pertama dan gabungin ke hasil query select kedua
	query := `SELECT p.name FROM permissions as p
			  JOIN role_has_permissions as rhp ON p.id = rhp.permission_id
			  JOIN user_has_roles as uhr ON rhp.role_id = uhr.role_id
			  WHERE uhr.user_id = ?
	
			  UNION

			  SELECT p.name FROM permissions as p
			  JOIN user_has_permissions as uhp ON p.id = uhp.permission_id
			  WHERE uhp.user_id = ?
			  `
	
	var binID []byte
	binID, err := userID.MarshalBinary()
	if err != nil {
		return  nil, err
	}

	rows, err := p.db.QueryContext(ctx, query,
		binID,
		binID,
	)
	if err != nil {
		return  nil, err
	}
	defer rows.Close()
	
	var permissions []string
	for rows.Next() {
		var permission string
		err := rows.Scan(
			&permission,
		)
		if err != nil {
			return  nil, err
		}
		permissions = append(permissions, permission)
	}

	// pastikan tidak ada error yg terjadi saat proses iterasi
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return  permissions, nil
}

// buat ambil permission berdasarkan id role
func (p *permissionRepository) GetPermissionsByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]string, error) {
	// 1. Guard clause: Jika tidak ada role, langsung kembalikan array kosong
	if len(roleIDs) == 0 {
		return []string{}, nil
	}

	// 2. Buat parameter dinamis (?, ?, ?) sesuai jumlah roleIDs
	placeholders := make([]string, len(roleIDs))
	args := make([]any, len(roleIDs))

	for i, id := range roleIDs {
		placeholders[i] = "?"
		args[i] = id[:] // Ubah tiap UUID menjadi []byte agar cocok dengan BINARY(16)
	}

	placeholderStr := strings.Join(placeholders, ",")

	// 3. Gunakan DISTINCT agar tidak ada nama permission yang duplikat
	// Gunakan IN (%s) untuk memasukkan tanda tanya secara dinamis
	query := fmt.Sprintf(`
		SELECT DISTINCT p.name 
		FROM permissions as p
		JOIN role_has_permissions as rhp ON p.id = rhp.permission_id
		WHERE rhp.role_id IN (%s)
	`, placeholderStr)

	// 4. Eksekusi query dengan menyebarkan args (...)
	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var permission string
		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}


