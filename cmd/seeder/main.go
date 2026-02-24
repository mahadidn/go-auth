package main

import (
	"log"
	"log/slog"

	"golang-auth/internal/config"
	"golang-auth/internal/seeder"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	slog.Info("Memulai Database Seeder...")

	// 1. Inisialisasi Database 
	db, err := config.NewDB() 
	if err != nil {
		slog.Error("Gagal terhubung ke database", "error", err)
		return // Hentikan seeder kalau DB mati
	}
	defer db.Close() // Pastikan koneksi ditutup setelah seeder selesai

	// 2. Jalankan fungsi Seeder
	seeder.SeedPermissionsAndSuperadmin(db)

	// Nanti kalau ada seeder lain, tinggal panggil di sini:
	// seeder.SeedRoles(db)
	// seeder.SeedAdminUser(db)

	slog.Info("Semua Seeder berhasil dieksekusi!")
}