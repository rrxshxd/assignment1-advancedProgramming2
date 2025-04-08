package usecase

import (
	"errors"
	"fmt"
	"github.com/rrxshxd/assignment1_advProg2/inventory_service/internal/entity"
	"github.com/rrxshxd/assignment1_advProg2/inventory_service/internal/repository"
	"time"
)

type ProductUseCase struct {
	productRepo repository.ProductRepository
}

func NewProductUseCase(productRepo repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{productRepo: productRepo}
}

func (uc *ProductUseCase) CreateProduct(product *entity.Product) error {
	if product.Name == "" || product.Category == "" || product.Category == "" {
		return errors.New("Invalid product data")
	}

	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now
	err := uc.productRepo.Create(product)
	if err != nil {
		return err
	}

	return nil
}

func (uc *ProductUseCase) GetProduct(id uint) (*entity.Product, error) {
	product, err := uc.productRepo.FindByID(id)
	if err != nil {
		if err == fmt.Errorf("Product with such id can't be found: %d", id) {
			return nil, fmt.Errorf("Product with such id can't be found: %d", id)
		}
		return nil, err
	}

	return product, nil
}

func (uc *ProductUseCase) UpdateProduct(product *entity.Product) error {
	existingProduct, err := uc.productRepo.FindByID(product.ID)
	if err != nil {
		if err == fmt.Errorf("Product with such id can't be found: %d", product.ID) {
			return fmt.Errorf("Product with such id can't be found: %d", product.ID)
		}
		return err
	}

	if product.Name != "" {
		existingProduct.Name = product.Name
	}
	if product.Description != "" {
		existingProduct.Description = product.Description
	}
	if product.Category != "" {
		existingProduct.Category = product.Category
	}
	if product.Price > 0 {
		existingProduct.Price = product.Price
	}
	if product.Stock >= 0 { // 0 is valid for stock
		existingProduct.Stock = product.Stock
	}

	existingProduct.UpdatedAt = time.Now()

	err = uc.productRepo.Update(existingProduct)
	if err != nil {
		return err
	}

	*product = *existingProduct

	return nil
}

func (uc *ProductUseCase) DeleteProduct(id uint) error {
	err := uc.productRepo.Delete(id)
	if err != nil {
		if err == fmt.Errorf("Product with such id can't be found: %d", id) {
			return fmt.Errorf("Product with such id can't be found: %d", id)
		}
		return err
	}

	return nil
}

func (uc *ProductUseCase) GetAll(page, limit int, filters map[string]interface{}) ([]*entity.Product, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	products, err := uc.productRepo.FindAll(page, limit, filters)
	if err != nil {
		return nil, err
	}

	return products, nil
}
