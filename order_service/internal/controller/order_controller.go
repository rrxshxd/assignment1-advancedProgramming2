package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rrxshxd/assignment1_advProg2/order_service/internal/entity"
	"github.com/rrxshxd/assignment1_advProg2/order_service/internal/usecase"
	"net/http"
	"strconv"
)

type OrderController struct {
	orderUseCase *usecase.OrderUseCase
}

func NewOrderController(orderUseCase *usecase.OrderUseCase) *OrderController {
	return &OrderController{orderUseCase: orderUseCase}
}

func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var order entity.Order
	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.orderUseCase.CreateOrder(&order); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, order)
}

func (c *OrderController) GetOrder(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	order, err := c.orderUseCase.GetOrder(uint(id))
	if err != nil {
		if err == errors.New("order not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (c *OrderController) UpdateOrderStatus(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	var request struct {
		Status entity.OrderStatus `json:"status" binding:"required,oneof=pending completed cancelled"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.orderUseCase.UpdateOrderStatus(uint(id), request.Status); err != nil {
		if err == errors.New("order not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *OrderController) GetUserOrders(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Query("user_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	orders, err := c.orderUseCase.GetUserOrders(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": orders})
}
