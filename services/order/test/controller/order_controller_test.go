package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	httpcontroller "product-app/services/order/internal/adapters/http/controller"
	"product-app/services/order/internal/domain"
	"product-app/services/order/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupOrderController() *httpcontroller.OrderController {
	initialOrders := []domain.Order{
		{
			Id:             1,
			CustomerNumber: "CUST-001",
			ProductID:      "PROD-1",
			Quantity:       2,
			OrderTime:      time.Now(),
		},
		{
			Id:             2,
			CustomerNumber: "CUST-002",
			ProductID:      "PROD-2",
			Quantity:       1,
			OrderTime:      time.Now(),
		},
	}

	fakeRepo := NewFakeOrderRepository(initialOrders)
	orderService := usecase.NewOrderService(fakeRepo)
	return httpcontroller.NewOrderController(orderService)
}

func Test_ShouldGetAllOrders(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := setupOrderController()

	err := controller.GetAllOrders(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var orders []map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &orders)
	assert.Equal(t, 2, len(orders))
}

func Test_ShouldGetOrderById(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	controller := setupOrderController()

	err := controller.GetOrderById(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, "CUST-001", response["customer_number"])
	assert.Equal(t, "PROD-1", response["product_id"])
}

func Test_ShouldGetOrdersByCustomerNumber(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/customer/CUST-001", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("customerNumber")
	c.SetParamValues("CUST-001")

	controller := setupOrderController()

	err := controller.GetByCustomerNumber(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var orders []map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &orders)
	assert.Equal(t, 1, len(orders))
	assert.Equal(t, "CUST-001", orders[0]["customer_number"])
}

func Test_ShouldCreateOrder(t *testing.T) {
	e := echo.New()
	payload := `{
		"customer_number": "CUST-003",
		"product_id": "PROD-9",
		"quantity": 3
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := setupOrderController()

	err := controller.CreateOrder(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func Test_ShouldCreateOrder_InvalidJSON(t *testing.T) {
	e := echo.New()
	payload := `{"customer_number": "CUST-003", "quantity": "bad"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := setupOrderController()

	err := controller.CreateOrder(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func Test_ShouldDeleteOrder(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/orders/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	controller := setupOrderController()

	err := controller.DeleteOrder(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
