package main

import (
	"context"

	"product-app/services/order/internal/adapters/http/controller"
	postgresql "product-app/services/order/internal/adapters/postgresql/common"
	"product-app/services/order/internal/config"
	"product-app/services/order/internal/usecase"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	ctx := context.Background()
	e := buildServer(ctx)
	e.Start(":8084")
}

func buildServer(ctx context.Context) *echo.Echo {
	e := echo.New()

	configurationManager := config.NewConfigurationManager()
	dbPool := postgresql.GetConnectionPool(ctx, configurationManager.PostgreSqlConfig)

	registerRoutes(e, dbPool)
	return e
}

func registerRoutes(e *echo.Echo, dbPool *pgxpool.Pool) {
	orderRepository := postgresql.NewOrderRepository(dbPool)
	orderService := usecase.NewOrderService(orderRepository)
	orderController := controller.NewOrderController(orderService)

	orderController.RegisterRoutes(e)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
}
