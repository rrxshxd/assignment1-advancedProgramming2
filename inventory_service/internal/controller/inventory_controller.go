package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rrxshxd/assignment1_advProg2/inventory_service/internal/entity"
	"github.com/rrxshxd/assignment1_advProg2/inventory_service/internal/usecase"
	"net/http"
	"strconv"
)

type InventoryController struct {
	productUseCase *usecase.ProductUseCase
}

func NewInventoryController(productUseCase *usecase.ProductUseCase) *InventoryController {
	return &InventoryController{productUseCase: productUseCase}
}

func (c *InventoryController) CreateProduct(ctx *gin.Context) {
	var request struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Category    string  `json:"category" binding:"required"`
		Price       float64 `json:"price" binding:"required,gt=0"`
		Stock       int     `json:"stock" binding:"gte=0"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := &entity.Product{
		Name:        request.Name,
		Description: request.Description,
		Category:    request.Category,
		Price:       request.Price,
		Stock:       request.Stock,
	}

	if err := c.productUseCase.CreateProduct(product); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := map[string]interface{}{
		"id":          product.ID,
		"name":        product.Name,
		"description": product.Description,
		"category":    product.Category,
		"price":       product.Price,
		"stock":       product.Stock,
		"created_at":  product.CreatedAt,
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *InventoryController) GetProduct(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	product, err := c.productUseCase.GetProduct(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	response := map[string]interface{}{
		"id":          product.ID,
		"name":        product.Name,
		"description": product.Description,
		"category":    product.Category,
		"price":       product.Price,
		"stock":       product.Stock,
		"created_at":  product.CreatedAt,
		"updated_at":  product.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *InventoryController) UpdateProduct(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	var request struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Category    string  `json:"category"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingProduct, err := c.productUseCase.GetProduct(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	if request.Name != "" {
		existingProduct.Name = request.Name
	}
	if request.Description != "" {
		existingProduct.Description = request.Description
	}
	if request.Category != "" {
		existingProduct.Category = request.Category
	}
	if request.Price > 0 {
		existingProduct.Price = request.Price
	}
	if request.Stock >= 0 {
		existingProduct.Stock = request.Stock
	}

	if err := c.productUseCase.UpdateProduct(existingProduct); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := map[string]interface{}{
		"id":          existingProduct.ID,
		"name":        existingProduct.Name,
		"description": existingProduct.Description,
		"category":    existingProduct.Category,
		"price":       existingProduct.Price,
		"stock":       existingProduct.Stock,
		"updated_at":  existingProduct.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *InventoryController) DeleteProduct(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	if err := c.productUseCase.DeleteProduct(uint(id)); err != nil {
		if err == fmt.Errorf("Product with such id can't be found: %d", id) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *InventoryController) GetAll(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	filters := make(map[string]interface{})
	if category := ctx.Query("category"); category != "" {
		filters["category"] = category
	}
	if minPrice := ctx.Query("min_price"); minPrice != "" {
		if min, err := strconv.ParseFloat(minPrice, 64); err == nil {
			filters["min_price"] = min
		}
	}
	if maxPrice := ctx.Query("max_price"); maxPrice != "" {
		if max, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filters["max_price"] = max
		}
	}
	if name := ctx.Query("name"); name != "" {
		filters["name"] = name
	}

	products, err := c.productUseCase.GetAll(page, limit, filters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]map[string]interface{}, 0, len(products))
	for _, product := range products {
		response = append(response, map[string]interface{}{
			"id":          product.ID,
			"name":        product.Name,
			"description": product.Description,
			"category":    product.Category,
			"price":       product.Price,
			"stock":       product.Stock,
			"created_at":  product.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  response,
		"page":  page,
		"limit": limit,
	})
}
