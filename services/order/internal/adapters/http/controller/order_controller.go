package controller

import "product-app/services/order/internal/usecase"

type OrderController struct {
	orderService usecase.IOrderService
}

func NewOrderController(orderService usecase.IOrderService) *OrderController {
	return &OrderController{orderService: orderService}
}
