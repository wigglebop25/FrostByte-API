package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"frostbyte-api/internal/domain"
	"frostbyte-api/internal/service"
	"github.com/go-chi/chi/v5"
)

type RoleHandler struct {
	service *service.RoleService
}

func NewRoleHandler(service *service.RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

func (h *RoleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var role domain.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateRole(&role); err != nil {
		if strings.Contains(err.Error(), "1062") {
			http.Error(w, "Role already exists", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(role)
}

func (h *RoleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	roles, err := h.service.GetAllRoles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(roles)
}

func (h *RoleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	role, err := h.service.GetRoleByID(uint(id))
	if err != nil {
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(role)
}

func (h *RoleHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var role domain.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	role.RoleID = uint(id)

	if err := h.service.UpdateRole(&role); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(role)
}

func (h *RoleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteRole(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *RoleHandler) AddPermission(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoleID     uint   `json:"role_id"`
		Permission string `json:"permission"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.AddPermission(req.RoleID, req.Permission); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *RoleHandler) AssignRole(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   uint   `json:"user_id"`
		RoleID   uint   `json:"role_id"`
		Username string `json:"username"`
		RoleName string `json:"role_name"` // e.g. "Customer" or "role_name" matching client payload
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Support both ID-based and Name-based assignment
	if req.Username != "" && req.RoleName != "" {
		if err := h.service.AssignRoleToUserByName(req.Username, req.RoleName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		if err := h.service.AssignRoleToUser(req.UserID, req.RoleID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	
	w.WriteHeader(http.StatusOK)
}
