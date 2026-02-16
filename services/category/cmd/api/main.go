package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"product-app/services/category/internal/adapters/http/controller"
	kafkaconsumer "product-app/services/category/internal/adapters/kafka"
	"product-app/services/category/internal/adapters/postgresql"
	pgcommon "product-app/services/category/internal/adapters/postgresql/common"
	"product-app/services/category/internal/config"
	"product-app/services/category/internal/usecase"
	sharedkafka "product-app/shared/kafka"

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

	go startKafkaConsumer()
}

func startKafkaConsumer() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	consumer := kafkaconsumer.NewConsumerAdapter(
		[]string{"kafka:9092"},
		"product.events",
		"category-service",
		func(ctx context.Context, message sharedkafka.Message) error {
			log.Printf("category-service received: %s", string(message.Value))
			return nil
		},
	)
	defer consumer.Close()

	if err := consumer.Start(ctx); err != nil {
		log.Printf("kafka consumer stopped: %v", err)
	}
}
