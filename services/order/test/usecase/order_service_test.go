package service

import (
	"testing"
	"time"

	"product-app/services/order/internal/domain"
	"product-app/services/order/internal/usecase"

	"github.com/stretchr/testify/assert"
)

func setupOrderService() usecase.IOrderService {
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
		{
			Id:             3,
			CustomerNumber: "CUST-001",
			ProductID:      "PROD-3",
			Quantity:       5,
			OrderTime:      time.Now(),
		},
	}

	fakeRepo := NewFakeOrderRepository(initialOrders)
	return usecase.NewOrderService(fakeRepo)
}

func Test_ShouldGetAllOrders(t *testing.T) {
	service := setupOrderService()
	orders, err := service.GetAll()
	assert.NoError(t, err)
	assert.Len(t, orders, 3)
}

func Test_ShouldGetOrderById(t *testing.T) {
	service := setupOrderService()
	order, err := service.GetById(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), order.Id)
	assert.Equal(t, "CUST-001", order.CustomerNumber)
}

func Test_ShouldGetOrdersByCustomerNumber(t *testing.T) {
	service := setupOrderService()
	orders, err := service.GetByCustomerNumber("CUST-001")
	assert.NoError(t, err)
	assert.Len(t, orders, 2)
}

func Test_ShouldCreateOrder(t *testing.T) {
	service := setupOrderService()
	before, _ := service.GetAll()
	assert.Len(t, before, 3)

	created, err := service.Create(domain.Order{
		CustomerNumber: "CUST-003",
		ProductID:      "PROD-9",
		Quantity:       3,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(4), created.Id)
	assert.Equal(t, "CUST-003", created.CustomerNumber)

	after, _ := service.GetAll()
	assert.Len(t, after, 4)
}

func Test_ShouldDeleteOrder(t *testing.T) {
	service := setupOrderService()
	err := service.Delete(2)
	assert.NoError(t, err)

	_, err = service.GetById(2)
	assert.Error(t, err)
}

func Test_ShouldFailValidation_WhenCustomerNumberMissing(t *testing.T) {
	service := setupOrderService()
	_, err := service.Create(domain.Order{
		CustomerNumber: "",
		ProductID:      "PROD-1",
		Quantity:       1,
	})
	assert.Error(t, err)
}

func Test_ShouldFailValidation_WhenQuantityInvalid(t *testing.T) {
	service := setupOrderService()
	_, err := service.Create(domain.Order{
		CustomerNumber: "CUST-009",
		ProductID:      "PROD-1",
		Quantity:       0,
	})
	assert.Error(t, err)
}
