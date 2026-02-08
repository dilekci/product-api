package infrastructure

import (
	"context"
	"os"
	"product-app/common/postgresql"
	"product-app/persistence"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	ctx                context.Context
	dbPool             *pgxpool.Pool
	productRepository  persistence.IProductRepository
	categoryRepository persistence.ICategoryRepository
	userRepository     persistence.IUserRepository
)

func TestMain(m *testing.M) {
	ctx = context.Background()

	createTestDatabase(ctx)

	dbPool = postgresql.GetConnectionPool(ctx, postgresql.Config{
		Host:     "localhost",
		Port:     "6432",
		DbName:   "productapp_unit_test",
		UserName: "postgres",
		Password: "postgres",
	})

	createSchema(ctx, dbPool)

	productRepository = persistence.NewProductRepository(dbPool)
	categoryRepository = persistence.NewCategoryRepository(dbPool)
	userRepository = persistence.NewUserRepository(dbPool)

	code := m.Run()

	dbPool.Close()
	os.Exit(code)
}

func createTestDatabase(ctx context.Context) {
	adminPool := postgresql.GetConnectionPool(ctx, postgresql.Config{
		Host:                  "localhost",
		Port:                  "6432",
		DbName:                "postgres",
		UserName:              "postgres",
		Password:              "postgres",
		MaxConnections:        1,
		MaxConnectionIdleTime: 5 * time.Second,
	})

	_, _ = adminPool.Exec(ctx, `DROP DATABASE IF EXISTS productapp_unit_test`)
	_, _ = adminPool.Exec(ctx, `CREATE DATABASE productapp_unit_test`)
	adminPool.Close()
}
func createSchema(ctx context.Context, pool *pgxpool.Pool) {
	_, err := pool.Exec(ctx, `
		DROP TABLE IF EXISTS product_images;
		DROP TABLE IF EXISTS products;
		DROP TABLE IF EXISTS categories;
		DROP TABLE IF EXISTS users;

		CREATE TABLE categories (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT
		);

		CREATE TABLE products (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			price REAL NOT NULL,
			description TEXT,
			discount REAL,
			store TEXT,
			category_id BIGINT
		);

		CREATE TABLE product_images (
			id BIGSERIAL PRIMARY KEY,
			product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
			image_url TEXT NOT NULL
		);

		CREATE TABLE users (
			id BIGSERIAL PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			first_name TEXT,
			last_name TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
	`)
	if err != nil {
		panic(err)
	}
}
