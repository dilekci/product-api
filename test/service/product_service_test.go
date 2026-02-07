package service

import (
	"os"
	"product-app/domain"
	"product-app/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

var productService service.IProductService

func TestMain(m *testing.M) {

	initialProducts := []domain.Product{
		{
			Id:    1,
			Name:  "AirFryer",
			Price: 1000.0,
			Store: "ABC TECH",
		},
		{
			Id:    2,
			Name:  "Blender",
			Price: 500.0,
			Store: "XYZ Appliances",
		},
	}

	fakeRepository := NewFakeProductRepository(initialProducts)

	productService = service.NewProductService(fakeRepository)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func Test_ShouldGetAllProduct(t *testing.T) {
	t.Run("ShouldGetAllProducts", func(t *testing.T) {
		actualProduct := productService.GetAllProducts()

		assert.Equal(t, 2, len(actualProduct))
	})
}

func Test_ShouldGetAllProductsByStore(t *testing.T) {
	t.Run("ShouldGetAllProductsByStore", func(t *testing.T) {
		actualStore := productService.GetAllProductsByStore("ABC TECH")

		assert.NotEmpty(t, actualStore)

		for _, product := range actualStore {
			assert.Equal(t, "ABC TECH", product.Store)
		}
	})
}

func Test_ShouldGetById(t *testing.T) {
	t.Run("ShouldGetById", func(t *testing.T) {

		product, err := productService.GetById(1)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), product.Id)
		assert.Equal(t, "ABC TECH", product.Store)
		assert.Equal(t, "AirFryer", product.Name)

	})
}

func Test_ShouldDeleteById(t *testing.T) {
	t.Run("ShouldDeleteById", func(t *testing.T) {

		err := productService.DeleteById(1)

		assert.NoError(t, err)

		_, err = productService.GetById(1)
		assert.Error(t, err)

	})
}

func Test_ShouldUpdatePrice(t *testing.T) {
	t.Run("ShouldUpdatePrice", func(t *testing.T) {

		before, _ := productService.GetById(1)
		assert.Equal(t, float32(1000), before.Price)

		err := productService.UpdatePrice(1, 4200)
		assert.NoError(t, err)

		after, _ := productService.GetById(1)
		assert.Equal(t, float32(4200), after.Price)
	})
}

func Test_ShouldDeleteAllProduct(t *testing.T) {
	t.Run("ShouldDeleteAllProduct", func(t *testing.T) {
		err := productService.DeleteAllProducts()
		assert.NoError(t, err)

		products := productService.GetAllProducts()
		assert.Len(t, products, 0)

	})
}
