package main

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	infrastructure2 "order-service/infrastructure/log"
	"order-service/internal/api"
	"order-service/internal/repository"
	"order-service/internal/service"
	middleware2 "order-service/middleware"
	"order-service/routes"
	"time"
)

func main() {
	infrastructure2.InitLogger()

	orderRepo := repository.NewOrderRepository()
	orderService := service.NewOrderService(orderRepo, "order-service/internal/config", "order-service/internal/config/config.yaml")
	orderHandler := api.NewOrderHandler(orderService)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiterWithConfig(middleware2.GetRateLimiter()))
	e.Use(middleware.ContextTimeout(15 * time.Second))
	e.Use(echojwt.JWT("secrete"))

	routes.SetupRoutes(e, orderHandler)
	e.Logger.Fatal(e.Start(":8080"))
}
