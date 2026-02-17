package main

import (
	"golang-auth/internal/config"
	"golang-auth/internal/pkg/logger"
	"log"
	"log/slog"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main()  {
	
	// load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// setup logger
	logger.SetupLogger()

	// inisialisasi DB
	db, err := config.NewDB()
	if err != nil {
		slog.Error("Gagal terkoneksi ke database")
		return
	}
	defer db.Close()

}
