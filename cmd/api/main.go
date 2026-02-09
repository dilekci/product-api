package main

import (
	"context"

	"product-app/internal/adapters/http/controller"
	"product-app/internal/adapters/postgresql"
	pgcommon "product-app/internal/adapters/postgresql/common"
	"product-app/internal/config/app"
	"product-app/internal/usecase"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

func main() {
	// Root context for application lifecycle
	ctx := context.Background()

	e := buildServer(ctx)

	// Start HTTP server
	e.Start("localhost:8080")
}

func buildServer(ctx context.Context) *echo.Echo {
	// Initialize Echo HTTP server
	e := echo.New()

	// Load application configuration
	configurationManager := app.NewConfigurationManager()

	// Initialize PostgreSQL connection pool
	dbPool := pgcommon.GetConnectionPool(ctx, configurationManager.PostgreSqlConfig)

	registerRoutes(e, dbPool)
	return e
}

func registerRoutes(e *echo.Echo, dbPool *pgxpool.Pool) {
	// --------------------
	// Product dependencies
	// --------------------
	productRepository := postgresql.NewProductRepository(dbPool)
	productService := usecase.NewProductService(productRepository)
	productController := controller.NewProductController(productService)

	// --------------------
	// Category dependencies
	// --------------------
	categoryRepository := postgresql.NewCategoryRepository(dbPool)
	categoryService := usecase.NewCategoryService(categoryRepository)
	categoryController := controller.NewCategoryController(categoryService)

	// --------------------
	// User dependencies
	// --------------------
	userRepository := postgresql.NewUserRepository(dbPool)
	userService := usecase.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	// Register HTTP routes
	productController.RegisterRoutes(e)
	categoryController.RegisterRoutes(e)
	userController.RegisterRoutes(e)
}
