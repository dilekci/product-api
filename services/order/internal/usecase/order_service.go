package usecase

import (
	"errors"
	"product-app/services/order/internal/domain"
	"product-app/services/order/internal/ports"
)

type IOrderService interface {
	GetAll() ([]domain.Order, error)
	GetById(id int64) (domain.Order, error)
	GetByCustomerNumber(customerNumber string) ([]domain.Order, error)
	Create(order domain.Order) (domain.Order, error)
	Delete(id int64) error
}

type OrderService struct {
	orderRepository ports.OrderRepository
}

func NewOrderService(orderRepository ports.OrderRepository) IOrderService {
	return &OrderService{
		orderRepository: orderRepository,
	}
}

// Create implements [IOrderService].
func (o *OrderService) Create(order domain.Order) (domain.Order, error) {
	if err := validateOrder(order); err != nil {
		return domain.Order{}, err
	}
	return o.orderRepository.Create(order)
}

// Delete implements [IOrderService].
func (o *OrderService) Delete(id int64) error {
	return o.orderRepository.Delete(id)
}

// GetAll implements [IOrderService].
func (o *OrderService) GetAll() ([]domain.Order, error) {
	return o.orderRepository.GetAll()
}

// GetByCustomerNumber implements [IOrderService].
func (o *OrderService) GetByCustomerNumber(customerNumber string) ([]domain.Order, error) {
	return o.orderRepository.GetByCustomerNumber(customerNumber)
}

// GetById implements [IOrderService].
func (o *OrderService) GetById(id int64) (domain.Order, error) {
	return o.orderRepository.GetById(id)
}

func validateOrder(order domain.Order) error {
	if order.CustomerNumber == "" {
		return errors.New("customer number is required")
	}
	if order.ProductID == "" {
		return errors.New("product id is required")
	}
	if order.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}
	return nil
}
