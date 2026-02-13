package infrastructure

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
)

func TruncateTestData(ctx context.Context, dbPool *pgxpool.Pool) {
	_, err := dbPool.Exec(ctx, `
		TRUNCATE TABLE categories RESTART IDENTITY CASCADE;
	`)
	if err != nil {
		log.Fatalf("❌ Failed to truncate test data: %v", err)
	}
	log.Infof("✅ Test tables truncated safely")
}
