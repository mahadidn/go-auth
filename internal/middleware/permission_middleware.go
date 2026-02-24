package middleware

import (
	"golang-auth/internal/domain"
	"golang-auth/internal/helper"
	"net/http"

	"github.com/google/uuid"
)

type PermissionMiddleware struct {
	permissionService domain.PermissionService
	roleService domain.RoleService
}

func NewPermissionMiddleware(ps domain.PermissionService, rs domain.RoleService) *PermissionMiddleware {
	return &PermissionMiddleware{
		permissionService: ps,
		roleService: rs,
	}
}

// Require adalah fungsi dinamis (wrapper) handler
func (m *PermissionMiddleware) Require(requirePerm string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Ambil userID dari context
		userID, ok := r.Context().Value(UserContextKey).(uuid.UUID)
		if !ok {
			helper.ResponseUnauthorized(w, "Sesi tidak valid atau tidak ditemukan")
			return
		}

		// ==========================================
		// TAHAP 1: Cek DIRECT PERMISSION terlebih dahulu
		// ==========================================
		directPermissions, err := m.permissionService.GetPermissionsByUserID(r.Context(), userID)
		if err != nil {
			helper.ResponseInternalError(w, "Gagal memverifikasi hak akses langsung")
			return
		}

		// Jika ketemu di direct permission, langsung beri akses dan BERHENTI (Early Return)
		for _, p := range directPermissions {
			if p == requirePerm {
				next(w, r)
				return // Hemat 2 query database!
			}
		}

		// ==========================================
		// TAHAP 2: Jika gagal, cari secara INDIRECT lewat ROLE
		// ==========================================
		roles, err := m.roleService.GetRoleByUserID(r.Context(), userID)
		if err != nil {
			helper.ResponseInternalError(w, "Gagal memverifikasi daftar role")
			return
		}

		// Guard clause: Kalau user ternyata nggak punya role sama sekali, langsung tolak
		if len(roles) == 0 {
			helper.ResponseForbidden(w, "Anda tidak memiliki izin untuk mengakses fitur ini")
			return
		}

		// Konversi roles string ke roles uuid
		var rolesUuid []uuid.UUID
		for _, role := range roles {
			roleUuid, err := uuid.Parse(role)
			// Tambahkan validasi error agar tidak memasukkan UUID kosong jika parsing gagal
			if err == nil {
				rolesUuid = append(rolesUuid, roleUuid)
			}
		}

		// Ambil permissions dari kumpulan role tersebut
		permissionsFromRoles, err := m.permissionService.GetPermissionsByRoleIDs(r.Context(), rolesUuid)
		if err != nil {
			helper.ResponseInternalError(w, "Gagal memverifikasi hak akses role")
			return
		}

		// Cek apakah requirePerm ada di daftar permission dari role
		for _, p := range permissionsFromRoles {
			if p == requirePerm {
				next(w, r)
				return // Beri akses dan berhenti
			}
		}

		// ==========================================
		// TAHAP 3: Tolak akses jika tidak ditemukan di mana-mana
		// ==========================================
		helper.ResponseForbidden(w, "Anda tidak memiliki izin untuk mengakses fitur ini")
	}
}