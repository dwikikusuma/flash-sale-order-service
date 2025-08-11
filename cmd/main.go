package main

import (
	"github.com/labstack/echo/v4"
	"order-service/internal/api"
	"order-service/internal/repository"
	"order-service/internal/service"
)

func main() {
	orderRepo := repository.NewOrderRepository()
	orderService := service.NewOrderService(orderRepo, "order-service/internal/config", "order-service/internal/config/config.yaml")
	orderHandler := api.NewOrderHandler(orderService)

	e := echo.New()
	orderHandler.RegisterRoutes(e)
	e.Logger.Fatal(e.Start(":8080"))
}
