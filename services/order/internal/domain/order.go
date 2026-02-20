package domain

import "time"

type Order struct {
	Id             int64     `json:"id"`
	CustomerNumber string    `json:"customer_number"`
	ProductID      string    `json:"product_id"`
	Quantity       int32     `json:"quantity"`
	OrderTime      time.Time `json:"order_time"`
}
