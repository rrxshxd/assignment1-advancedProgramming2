package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/rrxshxd/assignment1_advProg2/order_service/internal/config"
	"github.com/rrxshxd/assignment1_advProg2/order_service/internal/controller"
	"github.com/rrxshxd/assignment1_advProg2/order_service/internal/repository/postgres"
	"github.com/rrxshxd/assignment1_advProg2/order_service/internal/usecase"
	"log"
)

func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	orderRepo := postgres.NewOrderRepository(db)
	orderUseCase := usecase.NewOrderUseCase(orderRepo)
	orderController := controller.NewOrderController(orderUseCase)

	router := gin.Default()

	router.POST("/orders", orderController.CreateOrder)
	router.GET("/orders/:id", orderController.GetOrder)
	router.PATCH("/orders/:id", orderController.UpdateOrderStatus)
	router.GET("/orders", orderController.GetUserOrders)

	router.Run(":" + cfg.Port)
}
