package seeder

import (
	"database/sql"
	"log/slog"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// SeedPermissionsAndSuperadmin memastikan permission wajib dan superadmin tersedia
func SeedPermissionsAndSuperadmin(db *sql.DB) {
	// ==========================================
	// TAHAP 1: SEED PERMISSIONS
	// ==========================================
	permissions := []string{
		"roles:view",
		"roles:manage",
		"permissions:view",
		// Tambahkan yang lain jika ada
	}

	for _, name := range permissions {
		id, _ := uuid.NewV7()
		query := `INSERT IGNORE INTO permissions (id, name) VALUES (?, ?)`
		_, err := db.Exec(query, id[:], name)
		if err != nil {
			slog.Error("Gagal insert permission", "permission", name, "error", err)
		}
	}

	// ==========================================
	// TAHAP 2: AMBIL ID PERMISSION ASLI DARI DB
	// ==========================================
	// Kita ambil id asli dari DB karena id hasil generate tadi mungkin diabaikan oleh INSERT IGNORE
	rows, err := db.Query(`SELECT id, name FROM permissions`)
	if err != nil {
		slog.Error("Gagal mengambil daftar permissions", "error", err)
		return
	}
	defer rows.Close()

	// Simpan di struct sederhana untuk memudahkan iterasi
	var allPermissions []struct {
		ID   []byte
		Name string
	}
	for rows.Next() {
		var p struct {
			ID   []byte
			Name string
		}
		if err := rows.Scan(&p.ID, &p.Name); err == nil {
			allPermissions = append(allPermissions, p)
		}
	}

	// ==========================================
	// TAHAP 3: SEED USER SUPERADMIN
	// ==========================================
	emailAdmin := "superadmin@example.com"
	usernameAdmin := "superadmin"
	passwordAdmin := "superadmin123"

	// Hash password sebelum masuk DB
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordAdmin), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Gagal hashing password superadmin", "error", err)
		return
	}

	adminID, _ := uuid.NewV7()
	userQuery := `INSERT IGNORE INTO users (id, username, email, password) VALUES (?, ?, ?, ?)`
	_, err = db.Exec(userQuery, adminID[:], usernameAdmin, emailAdmin, string(hashedPassword))
	if err != nil {
		slog.Error("Gagal insert superadmin", "error", err)
		return
	}

	// Ambil ID Superadmin yang asli dari DB (berjaga-jaga kalau dia sudah pernah di-seed sebelumnya)
	var validAdminID []byte
	err = db.QueryRow(`SELECT id FROM users WHERE email = ?`, emailAdmin).Scan(&validAdminID)
	if err != nil {
		slog.Error("Gagal mengambil ID superadmin", "error", err)
		return
	}

	// ==========================================
	// TAHAP 4: HUBUNGKAN SUPERADMIN DENGAN SEMUA PERMISSIONS
	// ==========================================
	for _, p := range allPermissions {
		assignQuery := `INSERT IGNORE INTO user_has_permissions (user_id, permission_id) VALUES (?, ?)`
		result, err := db.Exec(assignQuery, validAdminID, p.ID)
		if err != nil {
			slog.Error("Gagal assign permission ke superadmin", "permission", p.Name, "error", err)
			continue
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			slog.Info("Permission di-assign ke Superadmin", "permission", p.Name)
		}
	}

	slog.Info("Proses Seeder Permission & Superadmin Selesai dengan Sukses!")
}