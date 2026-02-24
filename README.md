# Go-RBAC REST API

A robust **Role-Based Access Control (RBAC)** REST API built as a deep-dive study into the Go (Golang) backend ecosystem. Following Clean Architecture principles, this project is designed to be lightweight, secure, and highly performant without relying on heavy third-party web frameworks.

## 🎯 Learning Objectives & Case Study
This project serves as a practical implementation of advanced backend patterns in Go:
- Building a secure, layered middleware system (Authentication & Authorization).
- Handling complex database relationships using pure SQL and `database/sql`.
- Implementing dynamic, user-friendly error mapping for low-level database constraints.
- Managing standalone CLI tools (Seeders) alongside the main HTTP server.

## 🚀 Features

- **Pure Go Routing**: Built entirely using the standard `net/http` library (utilizing Go 1.22+ routing features).
- **Hybrid RBAC System**: Flexible access control supporting both *Direct Permissions* (assigned to users) and *Indirect Permissions* (inherited via roles).
- **Clean Architecture**: Strict separation of concerns between Domain, Service, Repository, and Handler layers.
- **Layered Security**: Sequential middleware execution separating token validation (Auth) and route-specific permission checks.
- **UUID v7 Integration**: Utilizing time-ordered UUIDs for primary keys to optimize MySQL indexing performance.
- **Idempotent Seeders**: Safe, CLI-driven database seeding for Superadmin and default permissions without data duplication.
- **Daily Rolling Logs**: Custom structured JSON logging with automated daily file rotation.
- **Thread-Safe Operations**: Optimized database connection pooling and concurrency-safe logger.
- **Dynamic Validation**: Comprehensive request payload validation with translated, user-friendly error messages.

## Run Locally
1. Clone the project

```bash
  git clone https://github.com/mahadidn/go-auth.git
```
2. Go to the project directory

```bash
  cd go-auth/
```
3. Install dependencies
```bash
  go mod tidy
```
4. Copy the example environment file and configure your database credentials
```bash
  cp .env.example .env
```
5. Setup database in the .env file and Run the SQL scripts located in the migrations/ folder to create the necessary tables.
6. Run Seeders, The superadmin account is located in internal/seeder/permission_seeder.go.
```bash
  go run cmd/seeder/main.go
```
7. Start the API server
```bash
  go run cmd/api/main.go
```



## 📁 Project Structure

```text
.
├── cmd/
│   ├── api/
│   │   └── main.go          # Main application entry point (HTTP Server)
│   └── seeder/
│       └── main.go          # CLI entry point for database seeding
├── internal/
│   ├── config/              # Database & Environment configurations
│   ├── domain/              # Core business models & Interfaces (Contracts)
│   ├── handler/             # HTTP Transport layer (Request/Response formatting)
│   ├── helper/              # Shared utilities (Error Mapper, JSON Writers)
│   ├── middleware/          # HTTP Interceptors (Auth, RBAC, Recovery)
│   ├── service/             # Business Logic layer
│   ├── repository/          # Database layer (Raw SQL queries)
│   ├── seeder/              # Seeder Layer
│   └── pkg/
│       └── logger/          # Custom Daily Log Writer implementation
├── logs/                    # Generated application log files (.log)
├── migrations/              # SQL Migration files for database schema
├── .env.example             # Environment variables template
└── go.mod                   # Go module dependencies

