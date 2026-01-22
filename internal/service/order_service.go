package service

import (
	"errors"
	"frostbyte-api/internal/domain"
	"frostbyte-api/internal/repository"
	"frostbyte-api/internal/websocket"
	"time"
)

type OrderService struct {
	repo        *repository.OrderRepository
	productRepo *repository.ProductRepository
	userRepo    *repository.UserRepository
	hub         *websocket.Hub
}

func NewOrderService(repo *repository.OrderRepository, productRepo *repository.ProductRepository, userRepo *repository.UserRepository, hub *websocket.Hub) *OrderService {
	return &OrderService{
		repo:        repo,
		productRepo: productRepo,
		userRepo:    userRepo,
		hub:         hub,
	}
}

func (s *OrderService) CreateOrder(userID uint, items []domain.OrderProduct) (*domain.Order, error) {
	var total float64
	for i, item := range items {
		product, err := s.productRepo.FindByID(item.ProductID)
		if err != nil {
			return nil, err
		}
		items[i].UnitPrice = product.Price
		items[i].LineTotal = product.Price * float64(item.Quantity)
		total += items[i].LineTotal
	}

	order := &domain.Order{
		UserID:      userID,
		TotalAmount: total,
		Status:      "PENDING",
		Products:    items,
	}

	if err := s.repo.Create(order); err != nil {
		return nil, err
	}

	// Broadcast the new order
	// Target: The User who created it + All Staff (Admin/Cashier)
	s.hub.Broadcast(order, order.UserID, []string{"Admin", "Cashier"})

	return order, nil
}

func (s *OrderService) GetAllOrders(filter domain.OrderFilter) ([]domain.Order, error) {
	return s.repo.GetAll(filter)
}

func (s *OrderService) GetOrderByID(id uint) (*domain.Order, error) {
	return s.repo.FindByID(id)
}

func (s *OrderService) GetOrdersByUserID(userID uint) ([]domain.Order, error) {
	return s.repo.GetByUserID(userID)
}

func (s *OrderService) GetOrdersByUser(username string) ([]domain.Order, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByUserID(user.UserID)
}

func (s *OrderService) GetOrdersByRole(role string) ([]domain.Order, error) {
	return s.repo.GetByRole(role)
}

func (s *OrderService) UpdateOrderStatus(id uint, status string) (*domain.Order, error) {
	// Validate Status
	allowedStatuses := map[string]bool{
		"PENDING":   true,
		"ACCEPTED":  true,
		"COOKING":   true,
		"READY":     true,
		"COMPLETED": true,
		"CANCELLED": true,
	}
	if !allowedStatuses[status] {
		return nil, errors.New("invalid status: allowed values are PENDING, ACCEPTED, COOKING, READY, COMPLETED, CANCELLED")
	}

	// Check if order exists
	order, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if err := s.repo.UpdateStatus(id, status); err != nil {
		return nil, err
	}

	// Fetch the updated order with all details to send in the broadcast
	updatedOrder, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Broadcast the full updated order
	// Target: The User who owns the order + All Staff (Admin/Cashier)
	s.hub.Broadcast(map[string]interface{}{
		"event": "ORDER_UPDATED",
		"data": map[string]interface{}{
			"order_id":   updatedOrder.OrderID,
			"status":     updatedOrder.Status,
			"updated_at": updatedOrder.UpdatedAt,
			"order":      updatedOrder, // Include full object for potential client use
		},
	}, order.UserID, []string{"Admin", "Cashier"})

	return updatedOrder, nil
}

func (s *OrderService) GetSalesAnalytics() (map[string]interface{}, error) {
	totalRevenue, totalOrders, statusCounts, dailyStats, err := s.repo.GetAnalytics()
	if err != nil {
		return nil, err
	}

	// Gap Filling Logic for the last 7 days
	// Create a map for quick lookup
	statsMap := make(map[string]repository.DailyRevenueStats)
	for _, stat := range dailyStats {
		// Ensure date format matches YYYY-MM-DD
		dateKey := stat.Date
		if len(dateKey) > 10 {
			dateKey = dateKey[:10]
		}
		statsMap[dateKey] = stat
		// DEBUG LOG
		println("DEBUG: DB Date:", stat.Date, " Key:", dateKey, " Rev:", stat.Revenue)
	}
	
	now := time.Now()
	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		println("DEBUG: Loop Date:", date) // DEBUG LOG
		// ... existing code ...

	var revenueTrend []repository.DailyRevenueStats
	// Use local time for generating the last 7 days.
	now := time.Now()
	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")

		if stat, exists := statsMap[date]; exists {
			revenueTrend = append(revenueTrend, stat)
		} else {
			revenueTrend = append(revenueTrend, repository.DailyRevenueStats{
				Date:       date,
				Revenue:    0,
				OrderCount: 0,
			})
		}
	}

	return map[string]interface{}{
		"summary": map[string]interface{}{
			"total_revenue": totalRevenue,
			"total_orders":  totalOrders,
			"status_counts": statusCounts,
		},
		"revenue_trend": revenueTrend,
	}, nil
}

func (s *OrderService) GetRevenueAnalytics(startDateStr, endDateStr string) (*domain.RevenueAnalyticsResponse, error) {
	// 1. Parse dates
	layout := "2006-01-02"
	start, err := time.Parse(layout, startDateStr)
	if err != nil {
		return nil, errors.New("invalid start_date format (YYYY-MM-DD)")
	}
	end, err := time.Parse(layout, endDateStr)
	if err != nil {
		return nil, errors.New("invalid end_date format (YYYY-MM-DD)")
	}

	if start.After(end) {
		return nil, errors.New("start_date cannot be after end_date")
	}

	// 2. Fetch data
	dailyStats, err := s.repo.GetRevenueAnalytics(startDateStr, endDateStr)
	if err != nil {
		return nil, err
	}

	// 3. Gap Filling
	statsMap := make(map[string]repository.DailyRevenueStats)
	for _, stat := range dailyStats {
		// Ensure date format is YYYY-MM-DD (take first 10 chars if it includes time)
		dateKey := stat.Date
		if len(dateKey) > 10 {
			dateKey = dateKey[:10]
		}
		statsMap[dateKey] = stat
	}

	var filledData []domain.DailyData
	var totalRevenue float64
	var totalOrders int64

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format(layout)
		stat, exists := statsMap[dateStr]

		revenue := 0.0
		orderCount := int64(0)

		if exists {
			revenue = stat.Revenue
			orderCount = stat.OrderCount
		}

		filledData = append(filledData, domain.DailyData{
			Date:       dateStr,
			DayName:    d.Weekday().String(),
			Revenue:    revenue,
			OrderCount: orderCount,
		})

		totalRevenue += revenue
		totalOrders += orderCount
	}

	return &domain.RevenueAnalyticsResponse{
		Period: domain.Period{
			Start: startDateStr,
			End:   endDateStr,
		},
		Summary: domain.Summary{
			TotalRevenue: totalRevenue,
			TotalOrders:  totalOrders,
		},
		DailyData: filledData,
	}, nil
}
