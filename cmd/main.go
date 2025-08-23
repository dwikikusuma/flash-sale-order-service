package main

import (
	"order-service/config"
	infrastructure "order-service/infrastructure/log"
	"order-service/internal/api"
	"order-service/internal/repository"
	"order-service/internal/service"
	reqMiddleware "order-service/middleware"
	"order-service/routes"
	"time"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	infrastructure.InitLogger()

	appConfig := config.LoadConfig(
		config.WithConfigFolder([]string{"./files/config"}),
		config.WithConfigFile("./files/config"),
		config.WithConfigType("yaml"),
	)

	orderRepo := repository.NewOrderRepository()
	orderService := service.NewOrderService(orderRepo, appConfig.Services.Product, appConfig.Services.Pricing)
	orderHandler := api.NewOrderHandler(orderService)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiterWithConfig(reqMiddleware.GetRateLimiter()))
	e.Use(middleware.ContextTimeout(15 * time.Second))
	e.Use(echojwt.JWT(appConfig.Secret.JWTSecret))

	routes.SetupRoutes(e, orderHandler)
	e.Logger.Fatal(e.Start(":" + appConfig.App.Port))
}
