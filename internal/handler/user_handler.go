package handler

import (
	"encoding/json"
	"golang-auth/internal/domain"
	"golang-auth/internal/helper"
	"golang-auth/internal/middleware"
	"net/http"

	"github.com/google/uuid"
)

type UserHandler struct {
	userService domain.UserService
}

func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request){
	// tangkap JSON dari request
	registerReq := &domain.UserRegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(registerReq)
	if err != nil {
		helper.ResponseBadRequest(w, "Format JSON tidak valid")
		return
	}

	createUserReq := domain.UserCreateRequest{
		Username: registerReq.Username,
		Email: registerReq.Email,
		Password: registerReq.Password,
	}

	// service
	err = h.userService.Create(r.Context(), createUserReq)
	if err != nil {
		helper.ResponseBadRequest(w, err.Error())
		return
	}

	helper.ResponseCreated(w, "Registrasi berhasil, silakan login")
}

// untuk ambil informasi user yg login
func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request){

	userID, ok := r.Context().Value(middleware.UserContextKey).(uuid.UUID)
	if !ok {
        helper.ResponseUnauthorized(w, "Gagal mengambil identitas user")
        return
    }

	user, err := h.userService.FindByID(r.Context(), userID)
	if err != nil {
		helper.ResponseBadRequest(w, err.Error())
	}

	helper.ResponseOK(w, user)

}
