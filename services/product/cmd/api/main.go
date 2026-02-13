package main

import (
	"context"

	"product-app/services/product/internal/adapters/http/controller"
	"product-app/services/product/internal/adapters/postgresql"
	pgcommon "product-app/services/product/internal/adapters/postgresql/common"
	"product-app/services/product/internal/config"
	"product-app/services/product/internal/usecase"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	ctx := context.Background()
	e := buildServer(ctx)
	e.Start(":8081")
}

func buildServer(ctx context.Context) *echo.Echo {
	e := echo.New()

	configurationManager := config.NewConfigurationManager()
	dbPool := pgcommon.GetConnectionPool(ctx, configurationManager.PostgreSqlConfig)

	registerRoutes(e, dbPool)
	return e
}

func registerRoutes(e *echo.Echo, dbPool *pgxpool.Pool) {
	productRepository := postgresql.NewProductRepository(dbPool)
	productService := usecase.NewProductService(productRepository)
	productController := controller.NewProductController(productService)

	productController.RegisterRoutes(e)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
}
