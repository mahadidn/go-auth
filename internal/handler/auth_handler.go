package handler

import (
	"encoding/json"
	"golang-auth/internal/domain"
	"golang-auth/internal/helper"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userService  domain.UserService
	tokenService domain.PersonalAccessTokenService
}

// constructor
func NewAuthHandler(userService domain.UserService, tokenService domain.PersonalAccessTokenService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		tokenService: tokenService,
	}
}


// method login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	
	// set header di awal, agar semua response otomatis berformat JSON
	w.Header().Set("Content-Type", "application/json")

	// decode JSON request
	loginRequest := &domain.UserLoginRequest{}
	err := json.NewDecoder(r.Body).Decode(loginRequest)
	
	if err != nil {
		helper.ResponseBadRequest(w, "Format JSON tidak valid")
		return // return agar eksekusi berhenti
	}

	deviceName := r.UserAgent()
	if deviceName == ""{
		deviceName = "Unknown Device"
	}

	user, err := h.userService.FindByEmail(r.Context(), loginRequest.Email)
	if err != nil {
		helper.ResponseUnauthorized(w, "Email atau password salah")
		return
	}

	// cek password request dan password asli
	err = bcrypt.CompareHashAndPassword([]byte(loginRequest.Password), []byte(user.Password))
	if err != nil {
		helper.ResponseUnauthorized(w, "Email atau password salah")
		return
	}

	// kalau password benar generate tokennya
	token, expiresAt, err := h.tokenService.Create(r.Context(), domain.PersonalAccessTokenRequest{
		UserID: user.ID,
		TokenName: deviceName,
	})
	if err != nil {
		helper.ResponseInternalError(w, "Gagal membuat token sistem")
		return
	}

	loginResponse := domain.UserLoginResponse{
		Id: user.ID.String(),
		Username: user.Username,
		Email: user.Email,
		Token: token,
		ExpiresAt: expiresAt,
	}

	helper.ResponseOK(w, loginResponse)
}

// method logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Ambil dari header Authorization
	authHeader := r.Header.Get("Authorization")

	// cek apakah headernya kosong atau formatnya bukan "Bearer <token>"
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer "){
		helper.ResponseUnauthorized(w, "Token tidak ditemukan atau format salah")
		return
	}

	// Potong tulisan "Bearer " (pajang 7 karakter termasuk spasi)
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// lempar ke service untuk dihapus dari database
	err := h.tokenService.Delete(r.Context(), tokenString)
	if err != nil {
		helper.ResponseInternalError(w, "Gagal melakukan logout")
		return
	}

	helper.ResponseOK(w, "Logout berhasil")
}