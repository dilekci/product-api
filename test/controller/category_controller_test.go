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

func setupCategoryController() *httpcontroller.CategoryController {
	initialCategories := []domain.Category{
		{Id: 1, Name: "Electronics", Description: "Electronic items"},
		{Id: 2, Name: "Books", Description: "Books and magazines"},
	}

	fakeRepo := NewFakeCategoryRepository(initialCategories)
	categoryService := usecase.NewCategoryService(fakeRepo)
	return httpcontroller.NewCategoryController(categoryService)
}

func Test_ShouldGetAllCategories(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	categoryController := setupCategoryController()

	err := categoryController.GetAllCategories(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var categories []map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &categories)
	assert.Equal(t, 2, len(categories))
}

func Test_ShouldGetCategoryById(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/categories/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	categoryController := setupCategoryController()

	err := categoryController.GetCategoryById(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.Equal(t, "Electronics", response["name"])
	assert.Equal(t, "Electronic items", response["description"])
}

func Test_ShouldAddCategory(t *testing.T) {
	e := echo.New()
	categoryJSON := `{
		"name": "Clothing",
		"description": "Apparel and accessories"
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/categories", strings.NewReader(categoryJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	categoryController := setupCategoryController()

	err := categoryController.AddCategory(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func Test_ShouldUpdateCategory(t *testing.T) {
	e := echo.New()
	categoryJSON := `{
		"name": "Electronics and Gadgets",
		"description": "Electronic items and smart gadgets"
	}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/categories/1", strings.NewReader(categoryJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	categoryController := setupCategoryController()

	err := categoryController.UpdateCategory(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func Test_ShouldDeleteCategory(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/categories/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	categoryController := setupCategoryController()

	err := categoryController.DeleteCategoryById(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
