package main

import (
	"fmt"
	"golang-auth/internal/config"
	"golang-auth/internal/handler"
	"golang-auth/internal/middleware"
	"golang-auth/internal/pkg/logger"
	"golang-auth/internal/repository"
	"golang-auth/internal/service"
	"log"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
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

	// inisialisasi validator
	validate := NewValidator()

	// wiring repository
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewPersonalAccessTokenRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	roleRepo := repository.NewRoleRepository(db)

	// wiring service
	userService := service.NewUserService(userRepo, db, validate)
	tokenService := service.NewPersonalAccessTokenService(tokenRepo, db, validate)
	permissionService := service.NewPermissionService(permissionRepo, db)
	roleService := service.NewRoleService(roleRepo, db, validate)

	// wiring handler & middleware
	authHandler := handler.NewAuthHandler(userService, tokenService)
	userHandler := handler.NewUserHandler(userService)
	permissionHandler := handler.NewPermissionHandler(permissionService)
	roleHandler := handler.NewRoleHandler(roleService)

	authMiddleware := middleware.NewAuthMiddleware(tokenService)
	permMiddleware := middleware.NewPermissionMiddleware(permissionService, roleService)

	// routing mux utama (publik)
	mux := http.NewServeMux()
	// route publik
	// mux.HandleFunc("GET /api/v1/", authHandler.TesPing)
	mux.HandleFunc("POST /api/v1/register", userHandler.Register)
	mux.HandleFunc("POST /api/v1/login", authHandler.Login)

	// route terproteksi middleware
	Group(mux, "/api/v1/", authMiddleware.Authenticate, func(subMux *http.ServeMux) {
		subMux.HandleFunc("POST /logout", authHandler.Logout)
		subMux.HandleFunc("GET /user", userHandler.Profile)

		subMux.HandleFunc("PUT /user/{id}/roles", userHandler.AssignRole)
		subMux.HandleFunc("PUT /user/{id}/permissions", userHandler.AssignPermission)

		subMux.HandleFunc("GET /roles", permMiddleware.Require("roles:view", roleHandler.FindAll))
        subMux.HandleFunc("GET /roles/{id}", permMiddleware.Require("roles:view", roleHandler.FindByID))
        subMux.HandleFunc("POST /roles", permMiddleware.Require("roles:manage", roleHandler.Create))
        subMux.HandleFunc("PUT /roles/{id}", permMiddleware.Require("roles:manage", roleHandler.Update))
        subMux.HandleFunc("DELETE /roles/{id}", permMiddleware.Require("roles:manage", roleHandler.Delete))

        subMux.HandleFunc("GET /permissions", permMiddleware.Require("permissions:view", permissionHandler.FindAll))
        subMux.HandleFunc("GET /permissions/user/{id}", permMiddleware.Require("permissions:view", permissionHandler.FindByUserID))
	})

	port := os.Getenv("APP_PORT")
	if port == ""{
		port = "8000"
	}

	fmt.Println("Server running", "port", port)

	handlerWithRecovery := middleware.Recovery(mux)

	server := &http.Server{
		Addr: ":" + port,
		Handler: handlerWithRecovery,
	}

	if err := server.ListenAndServe(); err != nil {
		slog.Error("Server stopped", "error", err)
	}

}

func Group(mux *http.ServeMux, prefix string, middleware func(http.Handler) http.Handler, route func(subMux *http.ServeMux)) {
    // buat mux kecil baru
	subMux := http.NewServeMux()

	// jalankan fungsi callback untuk mendaftarkan route didalam group
	route(subMux)

	// pasang submux ke mux utama dengan middleware
	mux.Handle(prefix, http.StripPrefix(strings.TrimSuffix(prefix, "/"), middleware(subMux)))
}

func NewValidator() *validator.Validate {
	v := validator.New()

	// Beritahu validator untuk memakai tag "json" sebagai nama field
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		// Ambil nilai dari tag json 
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		
		// Abaikan jika tag json-nya adalah "-"
		if name == "-" {
			return ""
		}
		
		return name
	})

	return v
}