package repository

import "github.com/rrxshxd/assignment1_advProg2/inventory_service/internal/entity"

type ProductRepository interface {
	Create(product *entity.Product) error
	FindByID(id uint) (*entity.Product, error)
	Update(product *entity.Product) error
	Delete(id uint) error
	FindAll(page, limit int, filters map[string]interface{}) ([]*entity.Product, error)
}
