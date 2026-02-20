package infrastructure

import (
	"testing"

	"product-app/services/order/internal/adapters/postgresql/common"
	"product-app/services/order/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestOrderRepository_GetAll(t *testing.T) {
	setupOrdersOnly()

	repo := postgresql.NewOrderRepository(dbPool)
	orders, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, orders, 3)
}

func TestOrderRepository_GetById(t *testing.T) {
	setupOrdersOnly()

	repo := postgresql.NewOrderRepository(dbPool)
	order, err := repo.GetById(1)
	assert.NoError(t, err)
	assert.Equal(t, "CUST-001", order.CustomerNumber)
}

func TestOrderRepository_GetByCustomerNumber(t *testing.T) {
	setupOrdersOnly()

	repo := postgresql.NewOrderRepository(dbPool)
	orders, err := repo.GetByCustomerNumber("CUST-001")
	assert.NoError(t, err)
	assert.Len(t, orders, 2)
}

func TestOrderRepository_Create(t *testing.T) {
	clearTestData()

	repo := postgresql.NewOrderRepository(dbPool)
	created, err := repo.Create(domain.Order{
		CustomerNumber: "CUST-009",
		ProductID:      "PROD-9",
		Quantity:       3,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), created.Id)
	assert.False(t, created.OrderTime.IsZero())
}

func TestOrderRepository_Delete(t *testing.T) {
	setupOrdersOnly()

	repo := postgresql.NewOrderRepository(dbPool)
	err := repo.Delete(1)
	assert.NoError(t, err)

	_, err = repo.GetById(1)
	assert.Error(t, err)
}
