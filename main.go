package main

import (
	"context"

	"product-app/common/app"
	"product-app/common/postgresql"
	"product-app/controller"
	"product-app/persistence"
	"product-app/service"

	"github.com/labstack/echo/v4"
)

func main() {
	// Root context for application lifecycle
	ctx := context.Background()

	// Initialize Echo HTTP server
	e := echo.New()

	// Load application configuration
	configurationManager := app.NewConfigurationManager()

	// Initialize PostgreSQL connection pool
	dbPool := postgresql.GetConnectionPool(ctx, configurationManager.PostgreSqlConfig)

	// --------------------
	// Product dependencies
	// --------------------
	productRepository := persistence.NewProductRepository(dbPool)
	productService := service.NewProductService(productRepository)
	productController := controller.NewProductController(productService)

	// --------------------
	// Category dependencies
	// --------------------
	categoryRepository := persistence.NewCategoryRepository(dbPool)
	categoryService := service.NewCategoryService(categoryRepository)
	categoryController := controller.NewCategoryController(categoryService)

	// --------------------
	// User dependencies
	// --------------------
	userRepository := persistence.NewUserRepository(dbPool)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	// Register HTTP routes
	productController.RegisterRoutes(e)
	categoryController.RegisterRoutes(e)
	userController.RegisterRoutes(e)

	// Start HTTP server
	e.Start("localhost:8080")
}
