package infrastructure

import (
	"context"
	"os"
	"product-app/services/order/internal/adapters/postgresql/common"
	"product-app/services/order/internal/ports"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	ctx             context.Context
	dbPool          *pgxpool.Pool
	orderRepository ports.OrderRepository
)

func TestMain(m *testing.M) {
	ctx = context.Background()

	host := getEnvString("TEST_DB_HOST", "localhost")
	port := getEnvString("TEST_DB_PORT", "6436")
	dbName := getEnvString("TEST_DB_NAME", "orderapp_unit_test")

	createTestDatabase(ctx, host, port, dbName)

	dbPool = postgresql.GetConnectionPool(ctx, postgresql.Config{
		Host:     host,
		Port:     port,
		DbName:   dbName,
		UserName: "postgres",
		Password: "postgres",
	})

	createSchema(ctx, dbPool)

	orderRepository = postgresql.NewOrderRepository(dbPool)
	code := m.Run()

	dbPool.Close()
	os.Exit(code)
}

func createTestDatabase(ctx context.Context, host, port, dbName string) {
	adminPool := postgresql.GetConnectionPool(ctx, postgresql.Config{
		Host:                  host,
		Port:                  port,
		DbName:                "postgres",
		UserName:              "postgres",
		Password:              "postgres",
		MaxConnections:        1,
		MaxConnectionIdleTime: 5 * time.Second,
	})

	_, _ = adminPool.Exec(ctx, `DROP DATABASE IF EXISTS `+dbName)
	_, _ = adminPool.Exec(ctx, `CREATE DATABASE `+dbName)
	adminPool.Close()
}

func createSchema(ctx context.Context, pool *pgxpool.Pool) {
	_, err := pool.Exec(ctx, `
		DROP TABLE IF EXISTS orders;

		CREATE TABLE orders (
			id BIGSERIAL PRIMARY KEY,
			customer_number TEXT NOT NULL,
			product_id TEXT NOT NULL,
			quantity INT NOT NULL CHECK (quantity > 0),
			order_time TIMESTAMP NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		panic(err)
	}
}

func getEnvString(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return parsed
}
