package middleware

import (
	"golang-auth/internal/helper"
	"log/slog"
	"net/http"
	"runtime/debug"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func () {
			if err := recover(); err != nil {
				// tangkap errornya
				slog.Error("CRITICAL PANIC RECOVERED",
					slog.Any("error", err),
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
				)

				// cetak stack trace
				slog.Error("Stack Trace: " + string(debug.Stack()))

				helper.ResponseInternalError(w, "Terjadi kesalahan internal pada server")
			}
		}()
		// lanjutkan request ke handler berikutnya
		next.ServeHTTP(w, r)
	})
}