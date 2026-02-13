package infrastructure

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
)

var insertCategories = `
INSERT INTO categories (name, description)
VALUES
('Elektronik', 'Elektronik ürünler'),
('Beyaz Eşya', 'Beyaz eşya ürünleri'),
('Dekorasyon', 'Ev dekorasyonu');
`

func InsertTestCategories(ctx context.Context, dbPool *pgxpool.Pool) {
	result, err := dbPool.Exec(ctx, insertCategories)
	if err != nil {
		log.Fatalf("❌ Failed to insert test categories: %v", err)
	}
	log.Infof("✅ Categories created: %d rows", result.RowsAffected())
}
