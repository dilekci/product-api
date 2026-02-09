package infrastructure

import (
	"testing"

	"product-app/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestProductRepository_GetAll(t *testing.T) {
	setupFullTestData()

	expected := []domain.Product{
		{Id: 1, Name: "AirFryer", Price: 3000, Description: "AirFryer açıklaması", Discount: 22, Store: "ABC TECH", CategoryID: 1},
		{Id: 2, Name: "Ütü", Price: 1500, Description: "Ütü açıklaması", Discount: 10, Store: "ABC TECH", CategoryID: 1},
		{Id: 3, Name: "Çamaşır Makinesi", Price: 10000, Description: "Çamaşır Makinesi açıklaması", Discount: 15, Store: "ABC TECH", CategoryID: 2},
		{Id: 4, Name: "Lambader", Price: 2000, Description: "Lambader açıklaması", Discount: 0, Store: "Dekorasyon Sarayı", CategoryID: 3},
	}

	actual := productRepository.GetAllProducts()

	assert.Len(t, actual, 4)
	assert.Equal(t, expected, actual)
}

func TestProductRepository_GetAllByStore(t *testing.T) {
	setupFullTestData()

	expected := []domain.Product{
		{Id: 1, Name: "AirFryer", Price: 3000, Description: "AirFryer açıklaması", Discount: 22, Store: "ABC TECH", CategoryID: 1},
		{Id: 2, Name: "Ütü", Price: 1500, Description: "Ütü açıklaması", Discount: 10, Store: "ABC TECH", CategoryID: 1},
		{Id: 3, Name: "Çamaşır Makinesi", Price: 10000, Description: "Çamaşır Makinesi açıklaması", Discount: 15, Store: "ABC TECH", CategoryID: 2},
	}

	actual := productRepository.GetAllProductsByStore("ABC TECH")

	assert.Equal(t, expected, actual)
}

func TestProductRepository_Add(t *testing.T) {
	clearTestData()

	newProduct := domain.Product{
		Name:        "Phone",
		Price:       3000,
		Description: "Apple phone",
		Discount:    0,
		Store:       "Apple Store",
		CategoryID:  0,
	}

	err := productRepository.AddProduct(newProduct)
	assert.NoError(t, err)

	products := productRepository.GetAllProducts()
	assert.Len(t, products, 1)
	assert.Equal(t, "Phone", products[0].Name)
}

func TestProductRepository_GetById(t *testing.T) {
	setupFullTestData()

	product, err := productRepository.GetById(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), product.Id)

	_, err = productRepository.GetById(999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product not found")
}

func TestProductRepository_DeleteById(t *testing.T) {
	setupFullTestData()

	err := productRepository.DeleteById(1)
	assert.NoError(t, err)

	_, err = productRepository.GetById(1)
	assert.Error(t, err)
}

func TestProductRepository_UpdatePrice(t *testing.T) {
	setupFullTestData()

	before, _ := productRepository.GetById(1)
	assert.Equal(t, float32(3000), before.Price)

	err := productRepository.UpdatePrice(1, 4000)
	assert.NoError(t, err)

	after, _ := productRepository.GetById(1)
	assert.Equal(t, float32(4000), after.Price)
}

func TestProductRepository_DeleteAll(t *testing.T) {
	clearTestData()

	err := productRepository.DeleteAllProducts()
	assert.NoError(t, err)

	products := productRepository.GetAllProducts()
	assert.Len(t, products, 0)
}
