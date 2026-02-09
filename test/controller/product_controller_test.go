package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	httpcontroller "product-app/internal/adapters/http/controller"
	"product-app/internal/domain"
	"product-app/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupProductController() *httpcontroller.ProductController {
	initialProducts := []domain.Product{
		{
			Id:          1,
			Name:        "AirFryer",
			Price:       1000,
			Description: "Digital air fryer",
			Discount:    10,
			Store:       "ABC TECH",
			CategoryID:  1,
		},
		{
			Id:          2,
			Name:        "Blender",
			Price:       500,
			Description: "High speed blender",
			Discount:    5,
			Store:       "XYZ Appliances",
			CategoryID:  1,
		},
	}

	fakeRepo := NewFakeProductRepository(initialProducts)
	productService := usecase.NewProductService(fakeRepo)
	return httpcontroller.NewProductController(productService)
}
func Test_ShouldGetProductId(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	productController := setupProductController()

	err := productController.GetProductById(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.Equal(t, "AirFryer", response["name"])
	assert.Equal(t, float64(1000), response["price"])
	assert.Equal(t, "Digital air fryer", response["description"])
	assert.Equal(t, float64(10), response["discount"])
	assert.Equal(t, "ABC TECH", response["store"])
	assert.Equal(t, float64(1), response["category_id"])
}

func Test_ShouldGetAllProducts(t *testing.T) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	productController := setupProductController()

	err := productController.GetAllProducts(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var products []map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &products)
	assert.Equal(t, 2, len(products))
}

func TestAddProduct_Success(t *testing.T) {
	e := echo.New()
	productJSON := `{
		"name": "Microwave",
		"price": 800,
		"description": "Digital microwave",
		"discount": 15,
		"store": "ABC TECH",
		"category_id": 1
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", strings.NewReader(productJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	productController := setupProductController()

	err := productController.AddProduct(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func Test_ShouldAddProduct_InvalidJSON(t *testing.T) {

	e := echo.New()
	invalidJSON := `{"name": "Test", "price": "invalid"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", strings.NewReader(invalidJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	productController := setupProductController()

	err := productController.AddProduct(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func Test_ShouldGetAllProducts_ByStore(t *testing.T) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/products?store=ABC+TECH", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	productController := setupProductController()

	err := productController.GetAllProducts(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var products []map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &products)
	assert.Equal(t, 1, len(products))
	assert.Equal(t, "AirFryer", products[0]["name"])
	assert.Equal(t, "ABC TECH", products[0]["store"])
}

func Test_ShouldAddProduct(t *testing.T) {
	e := echo.New()
	productJSON := `{
		"name": "Microwave",
		"price": 800,
		"description": "Digital microwave",
		"discount": 15,
		"store": "ABC TECH",
		"category_id": 1
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", strings.NewReader(productJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	productController := setupProductController()

	err := productController.AddProduct(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}
