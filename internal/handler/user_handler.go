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
		helper.ResponseBadRequest(w, helper.TranslateError(err))
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

// update
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	

}

// assign role
func (h *UserHandler) AssignRole(w http.ResponseWriter, r *http.Request) {
	
	userIDStr := r.PathValue("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		helper.ResponseBadRequest(w, "Format ID User tidak valid")
		return
	}

	assignRoleReq := &domain.AssignRoleRequest{}
	err = json.NewDecoder(r.Body).Decode(assignRoleReq)
	if err != nil {
		helper.ResponseBadRequest(w, "Format JSON tidak valid")
		return
	}

	err = h.userService.AssignRoles(r.Context(), userID, *assignRoleReq)
	if err != nil {
		helper.ResponseBadRequest(w, helper.TranslateError(err))
		return
	}

	helper.ResponseCreated(w, "Role pada user berhasil diubah")
}

func (h *UserHandler) AssignPermission(w http.ResponseWriter, r *http.Request) {
	
	userIDStr := r.PathValue("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		helper.ResponseBadRequest(w, "Format ID User tidak valid")
		return
	}

	assignPermReq := &domain.AssignPermissionRequest{}
	err = json.NewDecoder(r.Body).Decode(assignPermReq)
	if err != nil {
		helper.ResponseBadRequest(w, "Format JSON tidak valid")
		return
	}

	err = h.userService.AssignPermissions(r.Context(), userID, *assignPermReq)
	if err != nil {
		helper.ResponseBadRequest(w, helper.TranslateError(err))
		return
	}

	helper.ResponseCreated(w, "Permission pada user berhasil diubah")
}