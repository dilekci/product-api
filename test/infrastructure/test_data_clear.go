package infrastructure

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
)

const resetTablesQuery = `
DO $$
BEGIN
	IF to_regclass('public.product_images') IS NOT NULL THEN
		EXECUTE 'TRUNCATE TABLE product_images RESTART IDENTITY CASCADE';
	END IF;

	IF to_regclass('public.products') IS NOT NULL THEN
		EXECUTE 'TRUNCATE TABLE products RESTART IDENTITY CASCADE';
	END IF;

	IF to_regclass('public.categories') IS NOT NULL THEN
		EXECUTE 'TRUNCATE TABLE categories RESTART IDENTITY CASCADE';
	END IF;

	IF to_regclass('public.users') IS NOT NULL THEN
		EXECUTE 'TRUNCATE TABLE users RESTART IDENTITY CASCADE';
	END IF;
END $$;
`

func TruncateTestData(ctx context.Context, dbPool *pgxpool.Pool) {
	if _, err := dbPool.Exec(ctx, resetTablesQuery); err != nil {
		log.Fatalf("❌ Failed to truncate test tables: %v", err)
	}
	log.Info("✅ Test tables truncated safely")
}
