package helper

import (
	"encoding/json"
	"golang-auth/internal/domain"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, code int, status string, data any){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := domain.Response{
		Code: code,
		Status: status,
		Data: data,
	}

	_ = json.NewEncoder(w).Encode(response)
}

// Helper Khusus Sukses (200 OK)
func ResponseOK(w http.ResponseWriter, data any) {
    WriteJSON(w, http.StatusOK, "OK", data)
}

// Helper Khusus Created (201 Created)
func ResponseCreated(w http.ResponseWriter, data any) {
    WriteJSON(w, http.StatusCreated, "Created", data)
}

// Helper Khusus Client Error (400 Bad Request)
func ResponseBadRequest(w http.ResponseWriter, message string) {
    WriteJSON(w, http.StatusBadRequest, "Bad Request", message)
}

// Helper Khusus Auth Error (401 Unauthorized)
func ResponseUnauthorized(w http.ResponseWriter, message string) {
    WriteJSON(w, http.StatusUnauthorized, "Unauthorized", message)
}

// Helper Khusus Server Error (500 Internal Server Error)
func ResponseInternalError(w http.ResponseWriter, message string) {
    WriteJSON(w, http.StatusInternalServerError, "Internal Server Error", message)
}