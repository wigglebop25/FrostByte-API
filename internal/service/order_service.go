package service

import (
	"errors"
	"frostbyte-api/internal/domain"
	"frostbyte-api/internal/repository"
	"frostbyte-api/internal/websocket"
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
		total += product.Price * float64(item.Quantity)
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
	s.hub.Broadcast(order)

	return order, nil
}

func (s *OrderService) GetAllOrders() ([]domain.Order, error) {
	return s.repo.GetAll()
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

func (s *OrderService) UpdateOrderStatus(id uint, status string) error {
	// 1. Validate Status
	allowedStatuses := map[string]bool{
		"PENDING":   true,
		"ACCEPTED":  true,
		"COOKING":   true,
		"READY":     true,
		"COMPLETED": true,
		"CANCELLED": true,
	}
	if !allowedStatuses[status] {
		return errors.New("invalid status: allowed values are PENDING, ACCEPTED, COOKING, READY, COMPLETED, CANCELLED")
	}

	// 2. Check if order exists
	order, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("order not found")
	}

	if err := s.repo.UpdateStatus(id, status); err != nil {
		return err
	}

	// Broadcast the status update
	s.hub.Broadcast(map[string]interface{}{
		"order_id": id,
		"status":   status,
		"type":     "STATUS_UPDATE",
		"user_id":  order.UserID,
	})

	return nil
}

func (s *OrderService) GetSalesAnalytics() (map[string]interface{}, error) {
	totalRevenue, totalOrders, err := s.repo.GetAnalytics()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_revenue": totalRevenue,
		"total_orders":  totalOrders,
	}, nil
}
