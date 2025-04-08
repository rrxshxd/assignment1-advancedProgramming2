package repository

import "github.com/rrxshxd/assignment1_advProg2/order_service/internal/entity"

type OrderRepository interface {
	Create(order *entity.Order) error
	FindByID(id uint) (*entity.Order, error)
	UpdateStatus(id uint, status entity.OrderStatus) error
	FindByUserID(userID uint) ([]*entity.Order, error)
}
