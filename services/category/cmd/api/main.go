package main

import (
	"context"

	"product-app/services/category/internal/adapters/http/controller"
	"product-app/services/category/internal/adapters/postgresql"
	pgcommon "product-app/services/category/internal/adapters/postgresql/common"
	"product-app/services/category/internal/config"
	"product-app/services/category/internal/usecase"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	ctx := context.Background()
	e := buildServer(ctx)
	e.Start(":8082")
}

func buildServer(ctx context.Context) *echo.Echo {
	e := echo.New()

	configurationManager := config.NewConfigurationManager()
	dbPool := pgcommon.GetConnectionPool(ctx, configurationManager.PostgreSqlConfig)

	registerRoutes(e, dbPool)
	return e
}

func registerRoutes(e *echo.Echo, dbPool *pgxpool.Pool) {
	categoryRepository := postgresql.NewCategoryRepository(dbPool)
	categoryService := usecase.NewCategoryService(categoryRepository)
	categoryController := controller.NewCategoryController(categoryService)

	categoryController.RegisterRoutes(e)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
}
