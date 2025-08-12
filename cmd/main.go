package main

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"order-service/internal/api"
	"order-service/internal/infrastructure"
	"order-service/internal/repository"
	"order-service/internal/service"
)

func main() {
	infrastructure.InitLogger()

	orderRepo := repository.NewOrderRepository()
	orderService := service.NewOrderService(orderRepo, "order-service/internal/config", "order-service/internal/config/config.yaml")
	orderHandler := api.NewOrderHandler(orderService)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(echojwt.JWT("secrete"))

	orderHandler.RegisterRoutes(e)
	e.Logger.Fatal(e.Start(":8080"))
}
