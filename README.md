# Golang Auth

A backend REST API, following Clean Architecture principles. This project is designed to be lightweight, performant, and easy to maintain without relying on heavy third-party frameworks.

## 🚀 Features

- **Pure Go**: Built using `net/http`.
- **Clean Architecture**: Separation of concerns between Domain, Service, Repository, and Handler.
- **Daily Rolling Logs**: Custom structured logging (JSON) with daily file rotation.
- **Environment Management**: Configuration via `.env` files.
- **Thread-Safe**: Optimized database connection pooling and logger.

## 📁 Project Structure

```text
.
├── cmd/
│   ├── api/
│   │   └── main.go          # Entry point aplikasi
│   └── seeder/
│       └── main.go          # Untuk menjalankan seeder
├── internal/
│   ├── config/              # Konfigurasi Database & Environment
│   ├── domain/              # Model data & Interface (Kontrak)
│   ├── handler/             # Layer HTTP (Transport)
│   ├── helper/              # Helper
│   ├── middleware/          # Layer Middleware
│   ├── service/             # Layer Bisnis Logika
│   ├── repository/          # Layer Database (Raw SQL)
│   ├── seeder/              # Database Seeder
│   └── pkg/
│       └── logger/          # Custom Daily Log Writer
├── logs/                    # Folder penyimpanan log (.log)
├── migrations/              # SQL Migration files
├── .env.example             # Contoh konfigurasi env
└── go.mod                   # Module dependencies