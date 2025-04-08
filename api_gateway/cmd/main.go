package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rrxshxd/assignment1_advProg2/api_gateway/internal/config"
	"github.com/rrxshxd/assignment1_advProg2/api_gateway/internal/controller"
)

func main() {
	cfg := config.LoadConfig()

	gatewayController := controller.NewGatewayController(
		cfg.InventoryServiceURL,
		cfg.OrderServiceURL,
	)

	router := gin.Default()

	inventory := router.Group("/inventory")
	{
		inventory.GET("/products", gatewayController.ProxyInventory)
		inventory.POST("/products", gatewayController.ProxyInventory)
	}

	orders := router.Group("/orders")
	{
		orders.POST("/", gatewayController.ProxyOrders)
		orders.GET("/:id", gatewayController.ProxyOrders)
	}

	router.Run(":" + cfg.Port)
}
