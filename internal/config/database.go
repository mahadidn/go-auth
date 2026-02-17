package config

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

func NewDB() (*sql.DB, error) {

	// ambil dari env
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// susun DSN secara dinamis
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return  nil, err
	}

	// tes ping ke database
	err = db.Ping()
    if err != nil {
        return nil, err 
    }

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute) // 60 menit
	db.SetConnMaxIdleTime(10 * time.Minute) // 10 menit

	return db, nil
}
