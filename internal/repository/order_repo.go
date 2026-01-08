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
	return orders, err
}

func (r *OrderRepository) FindByID(id uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("Products.Product").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&domain.Order{}).Where("order_id = ?", id).Update("status", status).Error
}

func (r *OrderRepository) GetByUserID(userID uint) ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.Preload("Products.Product").Where("user_id = ?", userID).Find(&orders).Error
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
	return orders, err
}

func (r *OrderRepository) GetAnalytics() (float64, int64, error) {
	var totalRevenue float64
	var totalOrders int64

	if err := r.db.Model(&domain.Order{}).Select("COALESCE(SUM(total_amount), 0)").Scan(&totalRevenue).Error; err != nil {
		return 0, 0, err
	}

	if err := r.db.Model(&domain.Order{}).Count(&totalOrders).Error; err != nil {
		return 0, 0, err
	}

	return totalRevenue, totalOrders, nil
}
