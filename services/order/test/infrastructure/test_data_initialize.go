package infrastructure

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
)

var INSERT_ORDERS = `
INSERT INTO orders (customer_number, product_id, quantity, order_time)
VALUES
('CUST-001', 'PROD-1', 2, NOW()),
('CUST-002', 'PROD-2', 1, NOW()),
('CUST-001', 'PROD-3', 5, NOW());
`

func InsertTestOrders(ctx context.Context, dbPool *pgxpool.Pool) {
	result, err := dbPool.Exec(ctx, INSERT_ORDERS)
	if err != nil {
		log.Fatalf("❌ Failed to insert test orders: %v", err)
	}
	log.Infof("✅ Orders created: %d rows", result.RowsAffected())
}
