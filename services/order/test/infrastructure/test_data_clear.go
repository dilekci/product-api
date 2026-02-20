package infrastructure

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
)

const resetTablesQuery = `
DO $$
BEGIN
	IF to_regclass('public.orders') IS NOT NULL THEN
		EXECUTE 'TRUNCATE TABLE orders RESTART IDENTITY CASCADE';
	END IF;
END $$;
`

func TruncateTestData(ctx context.Context, dbPool *pgxpool.Pool) {
	if _, err := dbPool.Exec(ctx, resetTablesQuery); err != nil {
		log.Fatalf("❌ Failed to truncate test tables: %v", err)
	}
	log.Info("✅ Test tables truncated safely")
}
