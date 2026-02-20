package controller

import (
	"errors"
	"fmt"
	"time"

	"product-app/services/order/internal/domain"
	"product-app/services/order/internal/ports"
)

type FakeOrderRepository struct {
	orders []domain.Order
}

func NewFakeOrderRepository(initialOrders []domain.Order) ports.OrderRepository {
	return &FakeOrderRepository{orders: initialOrders}
}

func (repo *FakeOrderRepository) GetAll() ([]domain.Order, error) {
	return repo.orders, nil
}

func (repo *FakeOrderRepository) GetById(id int64) (domain.Order, error) {
	for _, order := range repo.orders {
		if order.Id == id {
			return order, nil
		}
	}
	return domain.Order{}, fmt.Errorf("order not found with id %d", id)
}

func (repo *FakeOrderRepository) GetByCustomerNumber(customerNumber string) ([]domain.Order, error) {
	var orders []domain.Order
	for _, order := range repo.orders {
		if order.CustomerNumber == customerNumber {
			orders = append(orders, order)
		}
	}
	if len(orders) == 0 {
		return nil, errors.New("no orders found for customer")
	}
	return orders, nil
}

func (repo *FakeOrderRepository) Create(order domain.Order) (domain.Order, error) {
	order.Id = int64(len(repo.orders)) + 1
	if order.OrderTime.IsZero() {
		order.OrderTime = time.Now()
	}
	repo.orders = append(repo.orders, order)
	return order, nil
}

func (repo *FakeOrderRepository) Delete(id int64) error {
	for i, order := range repo.orders {
		if order.Id == id {
			repo.orders = append(repo.orders[:i], repo.orders[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("order not found with id %d", id)
}
