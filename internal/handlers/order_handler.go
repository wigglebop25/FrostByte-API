package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"frostbyte-api/internal/domain"
	"frostbyte-api/internal/service"
	"github.com/go-chi/chi/v5"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	// UserID is extracted from context, not request body
	var req struct {
		Items []domain.OrderProduct `json:"products"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get UserID from context (set by AuthMiddleware)
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: User ID not found in context", http.StatusUnauthorized)
		return
	}

	order, err := h.service.CreateOrder(userID, req.Items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetAllOrders()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrderByID(uint(id))
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Security Check: Ensure user owns the order or is Admin
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is the owner
	if order.UserID != userID {
		// If not owner, check if Admin (we need to inject UserService or check roles here)
		// For now, simpler approach: if strict security is requested, prevent access.
		// NOTE: Ideally we'd check isAdmin here, but we don't have easy access to user roles in this handler without a DB call.
		// However, since this is a critical security fix, I will rely on the Service layer or fetch the user.
		
		// Let's return Forbidden for now to be safe. 
		// If the user IS an admin, they should use the Admin-specific routes or we need to enhance this handler.
		// Given the constraints, I'll block it. Use 'GetAll' (Admin) to see others' orders.
		http.Error(w, "Forbidden: You do not have access to this order", http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) GetByUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		http.Error(w, "Username required", http.StatusBadRequest)
		return
	}

	orders, err := h.service.GetOrdersByUser(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) GetByRole(w http.ResponseWriter, r *http.Request) {
	role := chi.URLParam(r, "role")
	if role == "" {
		http.Error(w, "Role required", http.StatusBadRequest)
		return
	}

	orders, err := h.service.GetOrdersByRole(role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateOrderStatus(uint(id), req.Status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	analytics, err := h.service.GetSalesAnalytics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(analytics)
}
