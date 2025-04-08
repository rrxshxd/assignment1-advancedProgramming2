package usecase

import (
	"github.com/rrxshxd/assignment1_advProg2/inventory_service/internal/entity"
	"github.com/rrxshxd/assignment1_advProg2/inventory_service/internal/repository"
)

type ProductUseCase struct {
	productRepo repository.ProductRepository
}

func NewProductUseCase(productRepo repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{productRepo: productRepo}
}

func (uc *ProductUseCase) CreateProduct(product *entity.Product) error {
	return uc.productRepo.Create(product)
}

func (uc *ProductUseCase) GetProduct(id uint) (*entity.Product, error) {
	return uc.productRepo.FindByID(id)
}

func (uc *ProductUseCase) UpdateProduct(product *entity.Product) error {
	return uc.productRepo.Update(product)
}

func (uc *ProductUseCase) DeleteProduct(id uint) error {
	return uc.productRepo.Delete(id)
}

//TODO: fix GetAll function
