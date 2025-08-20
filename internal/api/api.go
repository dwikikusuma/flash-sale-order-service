package api

import (
	_ "github.com/golang-jwt/jwt/v5"
	_ "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4/middleware"
	"order-service/internal/entity"
	"order-service/internal/service"
	"strconv"
)

type OrderHandler interface {
	CreateOrder(c echo.Context) error
	UpdateOrder(c echo.Context) error
	CancelOrder(c echo.Context) error
}

type orderHandler struct {
	OrderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) OrderHandler {
	return &orderHandler{
		OrderService: orderService,
	}
}

func (oh *orderHandler) CreateOrder(c echo.Context) error {
	var request entity.Order
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid order data"})
	}

	order, err := oh.OrderService.CreateOrder(&request)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to create order"})
	}

	return c.JSON(201, order)
}

func (oh *orderHandler) UpdateOrder(c echo.Context) error {
	var request entity.Order
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid order data"})
	}

	order, err := oh.OrderService.UpdateOrder(&request)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to update order"})
	}

	return c.JSON(200, order)
}

func (oh *orderHandler) CancelOrder(c echo.Context) error {
	orderIdStr := c.Param("id")
	orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid order ID"})
	}

	order, err := oh.OrderService.CancelOrder(orderId)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to cancel order"})
	}

	return c.JSON(200, order)
}
