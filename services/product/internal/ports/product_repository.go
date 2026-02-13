package ports

import "product-app/services/product/internal/domain"

type ProductRepository interface {
	GetAllProducts() []domain.Product
	GetProductsByCategoryId(categoryId int64) ([]domain.Product, error)
	GetAllProductsByStore(storeName string) []domain.Product
	AddProduct(product domain.Product) error
	GetById(productId int64) (domain.Product, error)
	DeleteById(productId int64) error
	UpdatePrice(productId int64, newPrice float32) error
	DeleteAllProducts() error
}
