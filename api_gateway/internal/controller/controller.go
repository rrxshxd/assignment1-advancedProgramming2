package controller

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

type GatewayController struct {
	inventoryServiceURL string
	orderServiceURL     string
}

func NewGatewayController(inventoryURL, orderURL string) *GatewayController {
	return &GatewayController{
		inventoryServiceURL: inventoryURL,
		orderServiceURL:     orderURL,
	}
}

func (c *GatewayController) ProxyInventory(ctx *gin.Context) {
	targetURL := c.inventoryServiceURL + ctx.Request.URL.Path

	req, err := http.NewRequest(ctx.Request.Method, targetURL, ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for key, values := range ctx.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Add X-Forwarded-For
	if clientIP := ctx.ClientIP(); clientIP != "" {
		req.Header.Add("X-Forwarded-For", clientIP)
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Copy response
	for key, values := range resp.Header {
		for _, value := range values {
			ctx.Header(key, value)
		}
	}
	ctx.Status(resp.StatusCode)
	ctx.Stream(func(w io.Writer) bool {
		io.Copy(w, resp.Body)
		return false
	})
}

func (c *GatewayController) ProxyOrders(ctx *gin.Context) {
	targetURL := c.orderServiceURL + ctx.Request.URL.Path

	req, err := http.NewRequest(ctx.Request.Method, targetURL, ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.Header = ctx.Request.Header.Clone()
	if authToken := ctx.GetHeader("Authorization"); authToken != "" {
		req.Header.Set("Authorization", authToken)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	ctx.DataFromReader(
		resp.StatusCode,
		resp.ContentLength,
		resp.Header.Get("Content-Type"),
		resp.Body,
		map[string]string{},
	)
}
