package usecase

import (
	"github.com/rrxshxd/assignment1_advProg2/order_service/internal/entity"
	"github.com/rrxshxd/assignment1_advProg2/order_service/internal/repository"
)

type OrderUseCase struct {
	orderRepo repository.OrderRepository
}

func NewOrderUseCase(orderRepo repository.OrderRepository) *OrderUseCase {
	return &OrderUseCase{orderRepo: orderRepo}
}

func (uc *OrderUseCase) CreateOrder(order *entity.Order) error {
	return uc.orderRepo.Create(order)
}

func (uc *OrderUseCase) GetOrder(id uint) (*entity.Order, error) {
	return uc.orderRepo.FindByID(id)
}

func (uc *OrderUseCase) UpdateOrderStatus(id uint, status entity.OrderStatus) error {
	return uc.orderRepo.UpdateStatus(id, status)
}

func (uc *OrderUseCase) GetUserOrders(userID uint) ([]*entity.Order, error) {
	return uc.orderRepo.FindByUserID(userID)
}
