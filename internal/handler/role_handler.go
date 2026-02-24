package handler

import (
	"encoding/json"
	"golang-auth/internal/domain"
	"golang-auth/internal/helper"
	"net/http"

	"github.com/google/uuid"
)

type RoleHandler struct {
	roleService domain.RoleService
}

func NewRoleHandler(roleService domain.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

func (h *RoleHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.roleService.FindAll(r.Context())
	if err != nil {
		helper.ResponseBadRequest(w, helper.TranslateError(err))
		return
	}

	helper.ResponseOK(w, data)
}

func (h *RoleHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	roleIDStr := r.PathValue("id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		helper.ResponseBadRequest(w, "Format ID Role tidak valid")
		return
	}
	data, err := h.roleService.FindById(r.Context(), roleID)
	if err != nil {
		helper.ResponseBadRequest(w, helper.TranslateError(err))
		return
	}

	helper.ResponseOK(w, data)

}

func (h *RoleHandler) Create(w http.ResponseWriter, r *http.Request){

	// tangkap JSON dari request
	createReq := &domain.RoleCreateRequest{}
	err := json.NewDecoder(r.Body).Decode(createReq)
	if err != nil {
		helper.ResponseBadRequest(w, "Format JSON tidak valid")
		return
	}
	
	err = h.roleService.Create(r.Context(), *createReq)
	if err != nil {
		helper.ResponseBadRequest(w, helper.TranslateError(err))
		return
	}

	helper.ResponseCreated(w, "Role berhasil ditambahkan")
}

func (h *RoleHandler) Update(w http.ResponseWriter, r *http.Request){
	// tangkap JSON dari request
	roleIDStr := r.PathValue("id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		helper.ResponseBadRequest(w, "Format ID Role tidak valid")
		return
	}

	updateRoleReq := &domain.RoleUpdateRequest{}
	err = json.NewDecoder(r.Body).Decode(updateRoleReq)
	if err != nil {
		helper.ResponseBadRequest(w, "Format JSON tidak valid")
		return
	}

	updateRoleReq.ID = roleID

	err = h.roleService.Update(r.Context(), *updateRoleReq)
	if err != nil {
		helper.ResponseBadRequest(w, helper.TranslateError(err))
		return
	}

	helper.ResponseOK(w, "Role berhasil diperbarui")
}

func (h *RoleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	
	roleIDStr := r.PathValue("id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		helper.ResponseBadRequest(w, "Format ID Role tidak valid")
		return
	}

	err = h.roleService.Delete(r.Context(), roleID)
	if err != nil {
		helper.ResponseBadRequest(w, helper.TranslateError(err))
		return
	}

	helper.ResponseOK(w, "Data berhasil terhapus")
}