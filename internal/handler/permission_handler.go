package handler

import (
	"golang-auth/internal/domain"
	"golang-auth/internal/helper"
	"net/http"

	"github.com/google/uuid"
)

type PermissionHandler struct {
	permissionService domain.PermissionService
}

func NewPermissionHandler(permissionService domain.PermissionService) *PermissionHandler {
	return &PermissionHandler{
		permissionService: permissionService,
	}
}

func (h *PermissionHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.permissionService.FindAll(r.Context())
	if err != nil {
		helper.ResponseBadRequest(w, helper.TranslateError(err))
		return
	}

	helper.ResponseOK(w, data)
}

func (h *PermissionHandler) FindByUserID(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		helper.ResponseBadRequest(w, "Format ID User tidak valid")
		return
	}

	data, err := h.permissionService.GetPermissionsByUserID(r.Context(), userID)
	if err != nil {
		helper.ResponseBadRequest(w, helper.TranslateError(err))
		return
	}

	helper.ResponseOK(w, data)

}