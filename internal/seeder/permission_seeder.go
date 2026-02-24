package seeder

import (
	"database/sql"
	"log/slog"

	"github.com/google/uuid"
)

// SeedPermissions berfungsi untuk memastikan permission wajib sudah ada di database
func SeedPermissions(db *sql.DB) {
	permissions := []string{
		"roles:view",
		"roles:manage",
		"permissions:view",
        // Nanti bisa tambah "users:view", "users:manage", dll kalau perlu
	}

	for _, name := range permissions {
		// Buat UUID v7 baru untuk ID
		id, err := uuid.NewV7()
		if err != nil {
			slog.Error("Gagal membuat UUID untuk seeder", "error", err)
			continue
		}

		// Gunakan INSERT IGNORE agar tidak error jika 'name' sudah ada
		query := `INSERT IGNORE INTO permissions (id, name) VALUES (?, ?)`

		// id[:] digunakan untuk mengubah [16]byte menjadi []byte agar sesuai dengan BINARY(16)
		result, err := db.Exec(query, id[:], name)
		if err != nil {
			slog.Error("Gagal menjalankan seeder permission", "permission", name, "error", err)
			continue
		}

		// Cek apakah data benar-benar dimasukkan atau dilewati karena sudah ada
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			slog.Info("Permission baru berhasil ditambahkan", "permission", name)
		}
	}
	
	slog.Info("Proses Seeder Permission Selesai!")
}