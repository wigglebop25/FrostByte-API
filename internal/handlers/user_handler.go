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

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		RoleName string `json:"role_name"` 
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.service.CreateUser(req.Username, req.Password, req.RoleName)
	if err != nil {
		if strings.Contains(err.Error(), "1062") {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}
		// Return specific role validation errors directly
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User added successfully"})
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Security Check
	tokenUserID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	roles, _ := GetRolesFromContext(r.Context())
	isAdmin := false
	for _, role := range roles {
		if role == "Admin" {
			isAdmin = true
			break
		}
	}

	// Allow if self-access or Admin
	if uint(id) != tokenUserID && !isAdmin {
		http.Error(w, "Forbidden: You do not have access to this user", http.StatusForbidden)
		return
	}

	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	// Security Check
	tokenUserID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.service.GetUserByID(tokenUserID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		RoleName string `json:"role_name"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := &domain.User{
		UserID:       uint(id),
		Username:     req.Username,
		PasswordHash: req.Password, // Service will hash this if not empty
	}

	if err := h.service.UpdateUser(user, req.RoleName); err != nil {
		if strings.Contains(err.Error(), "1062") {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Determine success message based on what changed
	var changes []string
	if req.Username != "" {
		changes = append(changes, "Username")
	}
	if req.Password != "" {
		changes = append(changes, "password")
	}
	if req.RoleName != "" {
		changes = append(changes, "role")
	}

	msg := "User updated successfully"
	if len(changes) == 1 {
		msg = "User " + changes[0] + " updated successfully"
	}

	json.NewEncoder(w).Encode(map[string]string{"message": msg})
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteUser(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("username")
	if query == "" {
		http.Error(w, "Query parameter 'username' required", http.StatusBadRequest)
		return
	}

	users, err := h.service.SearchUsers(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}
