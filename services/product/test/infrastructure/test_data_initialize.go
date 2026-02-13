package infrastructure

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
)

/* =========================
   CATEGORY TEST DATA
========================= */

var INSERT_CATEGORIES = `
INSERT INTO categories (name, description)
VALUES
('Elektronik', 'Elektronik ürünler'),
('Beyaz Eşya', 'Beyaz eşya ürünleri'),
('Dekorasyon', 'Ev dekorasyonu');
`

func InsertTestCategories(ctx context.Context, dbPool *pgxpool.Pool) {
	result, err := dbPool.Exec(ctx, INSERT_CATEGORIES)
	if err != nil {
		log.Fatalf("❌ Failed to insert test categories: %v", err)
	}
	log.Infof("✅ Categories created: %d rows", result.RowsAffected())
}

/* =========================
   PRODUCT TEST DATA
========================= */

var INSERT_PRODUCTS = `
INSERT INTO products (name, price, description, discount, store, category_id)
VALUES
('AirFryer', 3000.0, 'AirFryer açıklaması', 22.0, 'ABC TECH', 1),
('Ütü', 1500.0, 'Ütü açıklaması', 10.0, 'ABC TECH', 1),
('Çamaşır Makinesi', 10000.0, 'Çamaşır Makinesi açıklaması', 15.0, 'ABC TECH', 2),
('Lambader', 2000.0, 'Lambader açıklaması', 0.0, 'Dekorasyon Sarayı', 3);
`

func InsertTestProducts(ctx context.Context, dbPool *pgxpool.Pool) {
	result, err := dbPool.Exec(ctx, INSERT_PRODUCTS)
	if err != nil {
		log.Fatalf("❌ Failed to insert test products: %v", err)
	}
	log.Infof("✅ Products created: %d rows", result.RowsAffected())
}
