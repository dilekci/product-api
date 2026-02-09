package controller

import (
	"net/http"
	"product-app/internal/adapters/http/controller/response"
	"product-app/internal/domain"
	"product-app/internal/usecase"

	"github.com/labstack/echo/v4"
)

type CategoryController struct {
	categoryService usecase.ICategoryService
}

func NewCategoryController(categoryService usecase.ICategoryService) *CategoryController {
	return &CategoryController{categoryService: categoryService}
}

func (categoryController *CategoryController) RegisterRoutes(e *echo.Echo) {
	e.GET("/api/v1/categories", categoryController.GetAllCategories)
	e.GET("/api/v1/categories/:id", categoryController.GetCategoryById)
	e.POST("/api/v1/categories", categoryController.AddCategory)
	e.PUT("/api/v1/categories/:id", categoryController.UpdateCategory)
	e.DELETE("/api/v1/categories/:id", categoryController.DeleteCategoryById)
}

func (categoryController *CategoryController) GetAllCategories(c echo.Context) error {
	categories := categoryController.categoryService.GetAllCategories()
	return c.JSON(http.StatusOK, categories)
}

func (categoryController *CategoryController) GetCategoryById(c echo.Context) error {
	categoryId, err := parsePositiveIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid category ID",
		})
	}

	category, err := categoryController.categoryService.GetById(categoryId)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, category)
}

func (categoryController *CategoryController) AddCategory(c echo.Context) error {
	var category domain.Category
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	if err := categoryController.categoryService.AddCategory(category); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Category created successfully",
	})
}

func (categoryController *CategoryController) UpdateCategory(c echo.Context) error {
	categoryId, err := parsePositiveIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid category ID",
		})
	}

	var category domain.Category
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	category.Id = categoryId

	if err := categoryController.categoryService.UpdateCategory(category); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Category updated successfully",
	})
}

func (categoryController *CategoryController) DeleteCategoryById(c echo.Context) error {
	categoryId, err := parsePositiveIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid category ID",
		})
	}

	if err := categoryController.categoryService.DeleteById(categoryId); err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Category deleted successfully",
	})
}
