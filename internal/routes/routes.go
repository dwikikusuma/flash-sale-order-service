package routes

import (
	"github.com/labstack/echo/v4"
	"order-service/internal/api"
)

func SetupRoutes(e *echo.Echo, oh api.OrderHandler) {
	e.POST("/order", oh.CreateOrder)       // Create a new order
	e.PUT("/order", oh.UpdateOrder)        // Update an existing order
	e.DELETE("/order/:id", oh.CancelOrder) // Cancel an order by ID
}
