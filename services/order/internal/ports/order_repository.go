package ports

import "product-app/services/order/internal/domain"

type OrderRepository interface {
	GetAll() ([]domain.Order, error)
	GetById(id int64) (domain.Order, error)
	GetByCustomerNumber(customerNumber string) ([]domain.Order, error)
	Create(order domain.Order) (domain.Order, error)
	Delete(id int64) error
}
