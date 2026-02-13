package infrastructure

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"product-app/services/user/internal/adapters/postgresql"
	pgcommon "product-app/services/user/internal/adapters/postgresql/common"
	"product-app/services/user/internal/ports"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	ctx            context.Context
	dbPool         *pgxpool.Pool
	userRepository ports.UserRepository
)

func TestMain(m *testing.M) {
	ctx = context.Background()

	host := getEnvString("TEST_DB_HOST", "localhost")
	port := getEnvString("TEST_DB_PORT", "6435")
	dbName := getEnvString("TEST_DB_NAME", "userapp_unit_test")

	createTestDatabase(ctx, host, port, dbName)

	dbPool = pgcommon.GetConnectionPool(ctx, pgcommon.Config{
		Host:     host,
		Port:     port,
		DbName:   dbName,
		UserName: "postgres",
		Password: "postgres",
	})

	createSchema(ctx, dbPool)

	userRepository = postgresql.NewUserRepository(dbPool)

	code := m.Run()

	dbPool.Close()
	os.Exit(code)
}

func createTestDatabase(ctx context.Context, host, port, dbName string) {
	adminPool := pgcommon.GetConnectionPool(ctx, pgcommon.Config{
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
		DROP TABLE IF EXISTS users;

		CREATE TABLE users (
			id BIGSERIAL PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
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
