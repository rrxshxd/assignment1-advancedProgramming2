package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/rrxshxd/assignment1_advProg2/inventory_service/internal/config"
	"github.com/rrxshxd/assignment1_advProg2/inventory_service/internal/controller"
	"github.com/rrxshxd/assignment1_advProg2/inventory_service/internal/repository/postgres"
	"github.com/rrxshxd/assignment1_advProg2/inventory_service/internal/usecase"
	"log"
)

func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	productRepo := postgres.NewProductRepository(db)
	productUseCase := usecase.NewProductUseCase(productRepo)
	inventoryController := controller.NewInventoryController(productUseCase)

	router := gin.Default()

	router.POST("/products/create", inventoryController.CreateProduct)
	router.GET("/products/:id", inventoryController.GetProduct)
	router.PATCH("/products/:id", inventoryController.UpdateProduct)
	router.DELETE("/products/:id", inventoryController.DeleteProduct)
	router.GET("/products", inventoryController.GetAll)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
