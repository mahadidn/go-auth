package middleware

import (
	"context"
	"golang-auth/internal/domain"
	"golang-auth/internal/helper"
	"net/http"
	"strings"
)

type contextKey string

const UserContextKey contextKey = "user"
const TokenContextKey contextKey = "token"

type AuthMiddleware struct {
	tokenService domain.PersonalAccessTokenService
}

func NewAuthMiddleware(tokenService domain.PersonalAccessTokenService) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService: tokenService,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ambil header authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer "){
			helper.ResponseUnauthorized(w, "Token diperlukan")
			return 
		}
		// ekstrak token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// validasi token lewat service
		tokenData, err := m.tokenService.FindByToken(r.Context(), tokenString)
		if err != nil {
			helper.ResponseUnauthorized(w, "Token tidak valid atau sudah kedaluwarsa")
			return 
		}

		// masukkan userID kedalam context
		// supaya handler selanjutnya bisa tahu siapa yang sedang login
		ctx := context.WithValue(r.Context(), UserContextKey, tokenData.UserID)
		ctx = context.WithValue(ctx, TokenContextKey, tokenString)

		// update waktu penggunaan terakhir (menggunakan goroutine agar jalan secara asinkron)
		go m.tokenService.UpdateLastUsed(context.Background(), tokenString)

		// lanjutkan ke handler berikutnya
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}