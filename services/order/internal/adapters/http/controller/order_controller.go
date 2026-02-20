package controller

import (
	"net/http"
	"product-app/services/order/internal/adapters/http/controller/response"
	"product-app/services/order/internal/adapters/http/middleware"
	"product-app/services/order/internal/domain"
	"product-app/services/order/internal/usecase"

	"github.com/labstack/echo/v4"
)

type OrderController struct {
	orderService usecase.IOrderService
}

func NewOrderController(orderService usecase.IOrderService) *OrderController {
	return &OrderController{orderService: orderService}
}

func (orderController *OrderController) RegisterRoutes(e *echo.Echo) {
	e.GET("/api/v1/orders", orderController.GetAllOrders)
	e.GET("/api/v1/orders/:id", orderController.GetOrderById)
	e.GET("/api/v1/orders/customer/:customerNumber", orderController.GetByCustomerNumber)

	protected := e.Group("/api/v1/orders", middleware.JWTMiddleware())
	protected.POST("", orderController.CreateOrder)
	protected.DELETE("/:id", orderController.DeleteOrder)
}

func (orderController *OrderController) GetAllOrders(c echo.Context) error {
	orders, err := orderController.orderService.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, orders)
}

func (orderController *OrderController) GetOrderById(c echo.Context) error {
	orderId, err := parsePositiveIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid order ID",
		})
	}

	order, err := orderController.orderService.GetById(orderId)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, order)
}

func (orderController *OrderController) GetByCustomerNumber(c echo.Context) error {
	customerNumber := c.Param("customerNumber")
	if customerNumber == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid customer number",
		})
	}

	orders, err := orderController.orderService.GetByCustomerNumber(customerNumber)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, orders)
}

func (orderController *OrderController) CreateOrder(c echo.Context) error {
	var order domain.Order
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	created, err := orderController.orderService.Create(order)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, created)
}

func (orderController *OrderController) DeleteOrder(c echo.Context) error {
	orderId, err := parsePositiveIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid order ID",
		})
	}

	if err := orderController.orderService.Delete(orderId); err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Order deleted successfully",
	})
}
