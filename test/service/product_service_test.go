package service

import (
	"testing"

	"product-app/domain"
	"product-app/service"

	"github.com/stretchr/testify/assert"
)

func setupProductService() service.IProductService {
	initialProducts := []domain.Product{
		{
			Id:    1,
			Name:  "AirFryer",
			Price: 1000,
			Store: "ABC TECH",
		},
		{
			Id:    2,
			Name:  "Blender",
			Price: 500,
			Store: "XYZ Appliances",
		},
	}

	fakeRepository := NewFakeProductRepository(initialProducts)
	return service.NewProductService(fakeRepository)
}

func Test_ShouldGetAllProducts(t *testing.T) {
	productService := setupProductService()

	products := productService.GetAllProducts()

	assert.Len(t, products, 2)
}

func Test_ShouldGetAllProductsByStore(t *testing.T) {
	productService := setupProductService()

	products := productService.GetAllProductsByStore("ABC TECH")

	assert.NotEmpty(t, products)
	for _, product := range products {
		assert.Equal(t, "ABC TECH", product.Store)
	}
}

func Test_ShouldGetById(t *testing.T) {
	productService := setupProductService()

	product, err := productService.GetById(1)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), product.Id)
	assert.Equal(t, "AirFryer", product.Name)
	assert.Equal(t, "ABC TECH", product.Store)
}

func Test_ShouldDeleteById(t *testing.T) {
	productService := setupProductService()

	err := productService.DeleteById(1)
	assert.NoError(t, err)

	_, err = productService.GetById(1)
	assert.Error(t, err)
}

func Test_ShouldUpdatePrice(t *testing.T) {
	productService := setupProductService()

	before, err := productService.GetById(1)
	assert.NoError(t, err)
	assert.Equal(t, float32(1000), before.Price)

	err = productService.UpdatePrice(1, 4200)
	assert.NoError(t, err)

	after, err := productService.GetById(1)
	assert.NoError(t, err)
	assert.Equal(t, float32(4200), after.Price)
}

func Test_ShouldDeleteAllProducts(t *testing.T) {
	productService := setupProductService()

	err := productService.DeleteAllProducts()
	assert.NoError(t, err)

	products := productService.GetAllProducts()
	assert.Len(t, products, 0)
}
