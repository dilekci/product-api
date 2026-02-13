package controller

import (
	"errors"
	"fmt"
	"product-app/services/product/internal/domain"
	"product-app/services/product/internal/ports"
)

type FakeProductRepository struct {
	products []domain.Product
}

func NewFakeProductRepository(initialProducts []domain.Product) ports.ProductRepository {
	return &FakeProductRepository{
		products: initialProducts,
	}
}

func (fakeRepository *FakeProductRepository) GetAllProducts() []domain.Product {
	return fakeRepository.products
}

func (fakeRepository *FakeProductRepository) GetProductsByCategoryId(categoryId int64) ([]domain.Product, error) {
	var productsByCategory []domain.Product
	for _, product := range fakeRepository.products {
		if product.CategoryID == categoryId {
			productsByCategory = append(productsByCategory, product)
		}
	}
	if len(productsByCategory) == 0 {
		return nil, errors.New(fmt.Sprintf("No products found with category id %d", categoryId))
	}
	return productsByCategory, nil
}

func (fakeRepository *FakeProductRepository) GetAllProductsByStore(storeName string) []domain.Product {
	var productsByStore []domain.Product
	for _, product := range fakeRepository.products {
		if product.Store == storeName {
			productsByStore = append(productsByStore, product)
		}
	}
	return productsByStore
}

func (fakeRepository *FakeProductRepository) AddProduct(product domain.Product) error {
	fakeRepository.products = append(fakeRepository.products, domain.Product{
		Id:          int64(len(fakeRepository.products)) + 1,
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
		Discount:    product.Discount,
		Store:       product.Store,
		ImageUrls:   product.ImageUrls,
		CategoryID:  product.CategoryID,
	})
	return nil
}

func (fakeRepository *FakeProductRepository) GetById(productId int64) (domain.Product, error) {
	for _, product := range fakeRepository.products {
		if product.Id == productId {
			return product, nil
		}
	}
	return domain.Product{}, errors.New(fmt.Sprintf("Product not found with id %d", productId))
}

func (fakeRepository *FakeProductRepository) DeleteById(productId int64) error {
	foundIndex := -1
	for i, product := range fakeRepository.products {
		if product.Id == productId {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		return errors.New(fmt.Sprintf("Product not found with id %d", productId))
	}

	fakeRepository.products = append(fakeRepository.products[:foundIndex], fakeRepository.products[foundIndex+1:]...)
	return nil
}

func (fakeRepository *FakeProductRepository) UpdatePrice(productId int64, newPrice float32) error {
	found := false
	for i, product := range fakeRepository.products {
		if product.Id == productId {
			fakeRepository.products[i].Price = newPrice
			found = true
			break
		}
	}
	if !found {
		return errors.New(fmt.Sprintf("Product not found with id %d", productId))
	}
	return nil
}

func (fakeRepository *FakeProductRepository) DeleteAllProducts() error {
	fakeRepository.products = []domain.Product{}
	return nil
}
