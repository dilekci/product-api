package persistence

import (
	"context"
	"errors"
	"fmt"
	"product-app/domain"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
)

type IProductRepository interface {
	GettAllProducts() []domain.Product
	GetProductsByCategoryId(categoryId int64) ([]domain.Product, error)
	GetAllProductsByStore(storeName string) []domain.Product
	AddProduct(product domain.Product) error
	GetById(productId int64) (domain.Product, error)
	DeleteById(productId int64) error
	UpdatePrice(productId int64, newPrice float32) error
	DeleteAllProducts() error
}

type ProductRepository struct {
	dbPool *pgxpool.Pool
}

func NewProductRepository(dbPool *pgxpool.Pool) IProductRepository {
	return &ProductRepository{dbPool: dbPool}
}

func (r *ProductRepository) GettAllProducts() []domain.Product {
	ctx := context.Background()

	rows, err := r.dbPool.Query(ctx,
		`SELECT id, name, price, description, discount, store, category_id FROM products`)
	if err != nil {
		log.Errorf("❌ Error getting all products: %v", err)
		return []domain.Product{}
	}
	defer rows.Close()

	products, err := r.extractProducts(ctx, rows)
	if err != nil {
		log.Errorf("❌ Error extracting products: %v", err)
		return []domain.Product{}
	}
	return products
}

func (r *ProductRepository) GetAllProductsByStore(storeName string) []domain.Product {
	ctx := context.Background()

	rows, err := r.dbPool.Query(ctx, `
		SELECT id, name, price, description, discount, store, category_id
		FROM products
		WHERE store = $1
	`, storeName)

	if err != nil {
		log.Errorf("❌ Error querying products by store: %v", err)
		return []domain.Product{}
	}
	defer rows.Close()

	var products []domain.Product

	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.Price,
			&p.Description,
			&p.Discount,
			&p.Store,
			&p.CategoryID,
		); err != nil {
			log.Errorf("❌ Scan error: %v", err)
			continue
		}

		p.ImageUrls = r.loadImagesSafe(ctx, p.Id)
		products = append(products, p)
	}

	return products
}

func (r *ProductRepository) AddProduct(product domain.Product) error {
	ctx := context.Background()

	var productId int64
	err := r.dbPool.QueryRow(ctx, `
		INSERT INTO products (name, price, description, discount, store, category_id)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id
	`,
		product.Name,
		product.Price,
		product.Description,
		product.Discount,
		product.Store,
		product.CategoryID,
	).Scan(&productId)

	if err != nil {
		return fmt.Errorf("failed to insert product: %w", err)
	}

	if len(product.ImageUrls) == 0 {
		return nil
	}

	for _, url := range product.ImageUrls {
		_, err := r.dbPool.Exec(ctx,
			`INSERT INTO product_images (product_id, image_url) VALUES ($1,$2)`,
			productId, url,
		)
		if err != nil {
			log.Warnf("⚠️ Image insert failed for product %d: %v", productId, err)
		}
	}

	return nil
}

func (r *ProductRepository) GetById(productId int64) (domain.Product, error) {
	ctx := context.Background()

	var p domain.Product
	err := r.dbPool.QueryRow(ctx, `
		SELECT id, name, price, description, discount, store, category_id
		FROM products WHERE id = $1
	`, productId).Scan(
		&p.Id,
		&p.Name,
		&p.Price,
		&p.Description,
		&p.Discount,
		&p.Store,
		&p.CategoryID,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Product{}, fmt.Errorf("product not found with id %d", productId)
	}
	if err != nil {
		return domain.Product{}, err
	}

	p.ImageUrls = r.loadImagesSafe(ctx, productId)
	return p, nil
}

func (r *ProductRepository) DeleteById(productId int64) error {
	ctx := context.Background()

	ct, err := r.dbPool.Exec(ctx,
		`DELETE FROM products WHERE id = $1`, productId)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}
	log.Infof("✅ Product deleted with id %d", productId)
	return nil
}

func (r *ProductRepository) DeleteAllProducts() error {
	ctx := context.Background()

	_, err := r.dbPool.Exec(ctx, `DELETE FROM products`)
	if err != nil {
		return err
	}

	log.Info("✅ All products deleted (or already empty)")
	return nil
}

func (r *ProductRepository) UpdatePrice(productId int64, newPrice float32) error {
	ctx := context.Background()

	_, err := r.dbPool.Exec(ctx,
		`UPDATE products SET price = $1 WHERE id = $2`,
		newPrice, productId,
	)
	if err != nil {
		return err
	}

	log.Infof("✅ Product %d price updated to %v", productId, newPrice)
	return nil
}

func (r *ProductRepository) GetProductsByCategoryId(categoryId int64) ([]domain.Product, error) {
	ctx := context.Background()

	rows, err := r.dbPool.Query(ctx, `
		SELECT id, name, price, description, discount, store, category_id
		FROM products WHERE category_id = $1
	`, categoryId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product

	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.Price,
			&p.Description,
			&p.Discount,
			&p.Store,
			&p.CategoryID,
		); err != nil {
			return nil, err
		}

		p.ImageUrls = r.loadImagesSafe(ctx, p.Id)
		products = append(products, p)
	}

	return products, nil
}

func (r *ProductRepository) extractProducts(
	ctx context.Context,
	rows pgx.Rows,
) ([]domain.Product, error) {

	var products []domain.Product

	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.Price,
			&p.Description,
			&p.Discount,
			&p.Store,
			&p.CategoryID,
		); err != nil {
			return nil, err
		}

		p.ImageUrls = r.loadImagesSafe(ctx, p.Id)
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) loadImagesSafe(
	ctx context.Context,
	productId int64,
) []string {

	rows, err := r.dbPool.Query(ctx,
		`SELECT image_url FROM product_images WHERE product_id = $1`,
		productId,
	)
	if err != nil {
		log.Warnf("⚠️ Images not loaded for product %d: %v", productId, err)
		return []string{}
	}
	defer rows.Close()

	var images []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err == nil {
			images = append(images, url)
		}
	}
	return images
}
