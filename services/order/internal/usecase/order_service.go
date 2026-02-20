package usecase

import "product-app/services/order/internal/domain"

type IOrderService interface {
	GetAll() ([]domain.Order, error)
	GetById(id int64) (domain.Order, error)
	GetByCustomerNumber(customerNumber string) ([]domain.Order, error)
	Create(order domain.Order) (domain.Order, error)
	Delete(id int64) error
}

type OrderService struct {
}
