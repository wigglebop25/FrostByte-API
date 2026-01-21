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
	service     *service.OrderService
	userService *service.UserService
}

func NewOrderHandler(service *service.OrderService, userService *service.UserService) *OrderHandler {
	return &OrderHandler{
		service:     service,
		userService: userService,
	}
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
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Check if user is Admin or Cashier
	isStaff := false
	for _, role := range user.Roles {
		if role.Name == "Admin" || role.Name == "Cashier" {
			isStaff = true
			break
		}
	}

	var orders []domain.Order
	if isStaff {
		// Staff sees all orders
		orders, err = h.service.GetAllOrders()
	} else {
		// Customer only sees their own orders
		orders, err = h.service.GetOrdersByUserID(userID)
	}

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
		// Security Check: Ensure user owns the order.
		// Admins should use the Admin-specific routes.
		http.Error(w, "Forbidden: You do not have access to this order", http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) GetByUser(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "username")
	if param == "" {
		http.Error(w, "Username required", http.StatusBadRequest)
		return
	}

	// Security Check
	ctx := r.Context()
	tokenUsername, ok := GetUsernameFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	tokenUserID, ok := GetUserIDFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	roles, _ := GetRolesFromContext(ctx)
	isStaff := false
	for _, role := range roles {
		if role == "Admin" || role == "Cashier" {
			isStaff = true
			break
		}
	}

	var orders []domain.Order
	var err error

	// Check if the param is a User ID (numeric) or Username
	if targetID, errParse := strconv.ParseUint(param, 10, 32); errParse == nil {
		// It is a UserID
		if uint(targetID) != tokenUserID && !isStaff {
			http.Error(w, "Forbidden: You can only view your own orders", http.StatusForbidden)
			return
		}
		orders, err = h.service.GetOrdersByUserID(uint(targetID))
	} else {
		// It is a Username
		if param != tokenUsername && !isStaff {
			http.Error(w, "Forbidden: You can only view your own orders", http.StatusForbidden)
			return
		}
		orders, err = h.service.GetOrdersByUser(param)
	}

	if err != nil {
		// If user not found (or other error), return the specific message requested
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "No Results: No orders exist for the user " + param,
		})
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusOK) // Or NotFound? User said "No Results" message.
		json.NewEncoder(w).Encode(map[string]string{
			"message": "No Results: No orders exist for the user " + param,
		})
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
	json.NewEncoder(w).Encode(map[string]string{"message": "Status updated"})
}

func (h *OrderHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	analytics, err := h.service.GetSalesAnalytics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(analytics)
}

func (h *OrderHandler) GetRevenueAnalytics(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if startDate == "" || endDate == "" {
		http.Error(w, "start_date and end_date query parameters are required", http.StatusBadRequest)
		return
	}

	analytics, err := h.service.GetRevenueAnalytics(startDate, endDate)
	if err != nil {
		// Differentiate between validation errors and internal errors if needed, but error message usually helps
		if err.Error() == "invalid start_date format (YYYY-MM-DD)" || err.Error() == "invalid end_date format (YYYY-MM-DD)" || err.Error() == "start_date cannot be after end_date" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(analytics)
}
