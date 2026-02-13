package main

import (
	"context"

	"product-app/services/user/internal/adapters/http/controller"
	"product-app/services/user/internal/adapters/postgresql"
	pgcommon "product-app/services/user/internal/adapters/postgresql/common"
	"product-app/services/user/internal/config"
	"product-app/services/user/internal/usecase"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	ctx := context.Background()
	e := buildServer(ctx)
	e.Start(":8083")
}

func buildServer(ctx context.Context) *echo.Echo {
	e := echo.New()

	configurationManager := config.NewConfigurationManager()
	dbPool := pgcommon.GetConnectionPool(ctx, configurationManager.PostgreSqlConfig)

	registerRoutes(e, dbPool)
	return e
}

func registerRoutes(e *echo.Echo, dbPool *pgxpool.Pool) {
	userRepository := postgresql.NewUserRepository(dbPool)
	userService := usecase.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	userController.RegisterRoutes(e)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
}
