package repository

import (
	"frostbyte-api/internal/domain"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

type DailyRevenueStats struct {
	Date       string  `json:"date"`
	Revenue    float64 `json:"revenue"`
	OrderCount int64   `json:"order_count"`
}

func (r *OrderRepository) Create(order *domain.Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *OrderRepository) GetAll() ([]domain.Order, error) {
	var orders []domain.Order
	// Use Distinct to prevent duplicates from joins/preloads
	err := r.db.Distinct().Preload("Products.Product").Find(&orders).Error
	if err == nil {
		r.calculateLineTotals(orders)
	}
	return orders, err
}

func (r *OrderRepository) calculateLineTotals(orders []domain.Order) {
	for i := range orders {
		for j := range orders[i].Products {
			orders[i].Products[j].LineTotal = float64(orders[i].Products[j].Quantity) * orders[i].Products[j].UnitPrice
		}
	}
}

func (r *OrderRepository) FindByID(id uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("Products.Product").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	// Calculate line totals
	for j := range order.Products {
		order.Products[j].LineTotal = float64(order.Products[j].Quantity) * order.Products[j].UnitPrice
	}
	return &order, nil
}

func (r *OrderRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&domain.Order{}).Where("order_id = ?", id).Update("status", status).Error
}

func (r *OrderRepository) GetByUserID(userID uint) ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.Preload("Products.Product").Where("user_id = ?", userID).Find(&orders).Error
	if err == nil {
		r.calculateLineTotals(orders)
	}
	return orders, err
}

func (r *OrderRepository) GetByRole(roleName string) ([]domain.Order, error) {
	var orders []domain.Order
	// Join orders -> users -> user_roles -> roles
	err := r.db.Distinct("orders.*").Joins("User").
		Joins("JOIN user_roles ON user_roles.user_id = users.user_id").
		Joins("JOIN roles ON roles.role_id = user_roles.role_id").
		Where("roles.name = ?", roleName).
		Preload("Products.Product").
		Find(&orders).Error
	if err == nil {
		r.calculateLineTotals(orders)
	}
	return orders, err
}

func (r *OrderRepository) GetAnalytics() (float64, int64, map[string]int64, []DailyRevenueStats, error) {
	var totalRevenue float64
	var totalOrders int64

	// Calculate revenue only for COMPLETED and READY orders
	if err := r.db.Model(&domain.Order{}).Where("status IN ?", []string{"COMPLETED", "READY"}).Select("COALESCE(SUM(total_amount), 0)").Scan(&totalRevenue).Error; err != nil {
		return 0, 0, nil, nil, err
	}

	if err := r.db.Model(&domain.Order{}).Count(&totalOrders).Error; err != nil {
		return 0, 0, nil, nil, err
	}

	// Calculate counts for each status
	statusCounts := make(map[string]int64)
	rows, err := r.db.Model(&domain.Order{}).Select("status, count(*)").Group("status").Rows()
	if err != nil {
		return 0, 0, nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int64
		if err := rows.Scan(&status, &count); err != nil {
			return 0, 0, nil, nil, err
		}
		statusCounts[status] = count
	}

	// Ensure all standard statuses are present in the map, even if count is 0
	standardStatuses := []string{"PENDING", "ACCEPTED", "COOKING", "READY", "COMPLETED", "CANCELLED"}
	for _, status := range standardStatuses {
		if _, exists := statusCounts[status]; !exists {
			statusCounts[status] = 0
		}
	}

	// Get Daily Revenue Stats for the last 7 days
	var dailyStats []DailyRevenueStats
	// Using raw SQL for date grouping (MySQL compatible)
	err = r.db.Raw(`
		SELECT 
			DATE(created_at) as date, 
			SUM(total_amount) as revenue, 
			COUNT(order_id) as order_count 
		FROM orders 
		WHERE created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY) 
			AND status IN ('COMPLETED', 'READY') 
		GROUP BY DATE(created_at) 
		ORDER BY date ASC
	`).Scan(&dailyStats).Error

	if err != nil {
		return 0, 0, nil, nil, err
	}

	return totalRevenue, totalOrders, statusCounts, dailyStats, nil
}

func (r *OrderRepository) GetRevenueAnalytics(startDate, endDate string) ([]DailyRevenueStats, error) {
	var dailyStats []DailyRevenueStats
	err := r.db.Raw(`
		SELECT 
			DATE(created_at) as date, 
			SUM(total_amount) as revenue, 
			COUNT(order_id) as order_count 
		FROM orders 
		WHERE DATE(created_at) >= ? AND DATE(created_at) <= ? 
			AND status IN ('COMPLETED', 'READY') 
		GROUP BY DATE(created_at) 
		ORDER BY date ASC
	`, startDate, endDate).Scan(&dailyStats).Error

	if err != nil {
		return nil, err
	}
	return dailyStats, nil
}
