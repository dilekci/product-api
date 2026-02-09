package controller

import (
	"fmt"
	"net/http"
	"product-app/internal/adapters/http/controller/request"
	"product-app/internal/adapters/http/controller/response"
	"product-app/internal/adapters/http/middleware"
	"product-app/internal/usecase"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// ProductController handles HTTP requests for product operations
// It provides endpoints for CRUD operations on products with authentication support
type ProductController struct {
	productService usecase.IProductService
}

// NewProductController creates a new instance of ProductController
// Parameters:
//   - productService: Service interface for product business logic
//
// Returns:
//   - *ProductController: New controller instance
func NewProductController(productService usecase.IProductService) *ProductController {
	return &ProductController{productService: productService}
}

// RegisterRoutes registers all product-related HTTP routes
// Public routes (no authentication):
//   - GET /api/v1/products/:id - Get single product by ID
//   - GET /api/v1/products - Get all products (with optional store filter)
//
// Protected routes (JWT required):
//   - POST /api/v1/products - Create new product
//   - PUT /api/v1/products/:id - Update product price
//   - DELETE /api/v1/products/:id - Delete product by ID
//   - DELETE /api/v1/products/deleteAll - Delete all products
//   - GET /api/v1/products/my-products - Get current user's products
//
// Parameters:
//   - e: Echo instance for route registration
func (productController *ProductController) RegisterRoutes(e *echo.Echo) {
	// Public routes (no authentication required)
	e.GET("/api/v1/categories/:id/products", productController.GetProductsByCategoryId)
	e.GET("/api/v1/products/:id", productController.GetProductById)
	e.GET("/api/v1/products", productController.GetAllProducts)
	e.POST("/api/v1/products", productController.AddProduct)

	// Protected routes (authentication required)
	protected := e.Group("/api/v1/products", middleware.JWTMiddleware())
	protected.PUT("/:id", productController.UpdatePrice)
	protected.DELETE("/:id", productController.DeleteProductById)
	protected.DELETE("/deleteAll", productController.DeleteAllProducts)
}

func (productController *ProductController) GetProductsByCategoryId(c echo.Context) error {
	categoryId, err := parsePositiveIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid category ID",
		})
	}

	products, err := productController.productService.GetProductsByCategoryId(categoryId)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			Error: "Error: " + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, response.ToResponseList(products))
}

func (productController *ProductController) GetProductById(c echo.Context) error {
	productId, err := parsePositiveIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid product ID",
		})
	}

	product, err := productController.productService.GetById(productId)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			Error: "Error:  " + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, response.ToResponse(product))
}

func (productController *ProductController) GetAllProducts(c echo.Context) error {
	store := c.QueryParam("store")

	if len(store) == 0 {
		allProducts := productController.productService.GetAllProducts()
		return c.JSON(http.StatusOK, response.ToResponseList(allProducts))
	}
	productsWithGivenStore := productController.productService.GetAllProductsByStore(store)
	return c.JSON(http.StatusOK, response.ToResponseList(productsWithGivenStore))
}

func (productController *ProductController) AddProduct(c echo.Context) error {
	addProductRequest, bindErr := bindAddProductRequest(c)
	if bindErr != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: bindErr.Error(),
		})
	}
	err := productController.productService.Add(addProductRequest.ToModel())

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			Error: err.Error(),
		})
	}
	return c.NoContent(http.StatusCreated)
}
func (productController *ProductController) UpdatePrice(c echo.Context) error {
	productId, err := parsePositiveIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid product ID",
		})
	}

	newPrice, err := parsePriceQuery(c, "newPrice")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: err.Error(),
		})
	}

	productController.productService.UpdatePrice(productId, newPrice)
	return c.NoContent(http.StatusOK)
}

func (productController *ProductController) DeleteProductById(c echo.Context) error {
	productId, err := parsePositiveIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid product ID",
		})
	}
	err = productController.productService.DeleteById(productId)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			Error: err.Error(),
		})
	}
	return c.NoContent(http.StatusOK)
}

func (productController *ProductController) DeleteAllProducts(c echo.Context) error {
	err := productController.productService.DeleteAllProducts()
	if err != nil {
		log.Printf("DeleteAllProducts error: %v", err)
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			Error: err.Error(),
		})
	}
	return c.NoContent(http.StatusOK)
}

func bindAddProductRequest(c echo.Context) (request.AddProductRequest, error) {
	var addProductRequest request.AddProductRequest
	if err := c.Bind(&addProductRequest); err != nil {
		return request.AddProductRequest{}, err
	}
	return addProductRequest, nil
}

func parsePriceQuery(c echo.Context, name string) (float32, error) {
	raw := c.QueryParam(name)
	if len(raw) == 0 {
		return 0, fmt.Errorf("Parameter %s is required!", name)
	}
	convertedPrice, err := strconv.ParseFloat(raw, 32)
	if err != nil {
		return 0, fmt.Errorf("%s format disrupted!", name)
	}
	return float32(convertedPrice), nil
}
