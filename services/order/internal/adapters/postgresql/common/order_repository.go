package postgresql

import (
	"context"
	"fmt"
	"log"
	"product-app/services/order/internal/domain"
	"product-app/services/order/internal/ports"

	"github.com/jackc/pgx/v4/pgxpool"
)

type OrderRepository struct {
	dbPool *pgxpool.Pool
}

func NewOrderRepository(dbPool *pgxpool.Pool) ports.OrderRepository {
	return &OrderRepository{dbPool: dbPool}
}

// Create implements ports.OrderRepository
func (o *OrderRepository) Create(order domain.Order) (domain.Order, error) {
	ctx := context.Background()
	insertOrderSQL := `
		INSERT INTO orders (customer_number, product_id, quantity, order_time)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, order_time
	`
	err := o.dbPool.QueryRow(ctx, insertOrderSQL,
		order.CustomerNumber,
		order.ProductID,
		order.Quantity,
	).Scan(&order.Id, &order.OrderTime)
	if err != nil {
		log.Printf("❌ Error inserting order: %v", err)
		return domain.Order{}, fmt.Errorf("failed to insert order: %w", err)
	}
	log.Printf("✅ Order inserted with ID: %d", order.Id)
	return order, nil
}

// Delete implements ports.OrderRepository
func (o *OrderRepository) Delete(id int64) error {
	ctx := context.Background()
	deleteSql := `DELETE FROM orders WHERE id = $1`
	commandTag, err := o.dbPool.Exec(ctx, deleteSql, id)
	if err != nil {
		log.Printf("ERROR: Error while deleting order with id %d: %v", id, err)
		return fmt.Errorf("error while deleting order with id %d: %w", id, err)
	}
	if commandTag.RowsAffected() == 0 {
		log.Printf("WARNING: order with id %d not found for deletion", id)
		return fmt.Errorf("order with id %d not found", id)
	}
	log.Printf("INFO: order deleted with id %d", id)
	return nil
}

// GetAll implements ports.OrderRepository
func (o *OrderRepository) GetAll() ([]domain.Order, error) {
	ctx := context.Background()
	orderRows, err := o.dbPool.Query(ctx,
		`SELECT id, customer_number, product_id, quantity, order_time FROM orders`)
	if err != nil {
		return nil, fmt.Errorf("error while getting all orders: %w", err)
	}
	defer orderRows.Close()

	var orders []domain.Order
	for orderRows.Next() {
		var order domain.Order
		err := orderRows.Scan(&order.Id, &order.CustomerNumber, &order.ProductID, &order.Quantity, &order.OrderTime)
		if err != nil {
			log.Printf("ERROR: Error while scanning order: %v", err)
			continue
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// GetByCustomerNumber implements ports.OrderRepository
func (o *OrderRepository) GetByCustomerNumber(customerNumber string) ([]domain.Order, error) {
	ctx := context.Background()
	sql := `SELECT id, customer_number, product_id, quantity, order_time FROM orders WHERE customer_number = $1`
	rows, err := o.dbPool.Query(ctx, sql, customerNumber)
	if err != nil {
		return nil, fmt.Errorf("error while getting orders for customer %s: %w", customerNumber, err)
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		err := rows.Scan(&order.Id, &order.CustomerNumber, &order.ProductID, &order.Quantity, &order.OrderTime)
		if err != nil {
			log.Printf("ERROR: Error while scanning order: %v", err)
			continue
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// GetById implements ports.OrderRepository
func (o *OrderRepository) GetById(id int64) (domain.Order, error) {
	ctx := context.Background()
	sql := `SELECT id, customer_number, product_id, quantity, order_time FROM orders WHERE id = $1`
	var order domain.Order
	err := o.dbPool.QueryRow(ctx, sql, id).
		Scan(&order.Id, &order.CustomerNumber, &order.ProductID, &order.Quantity, &order.OrderTime)
	if err != nil {
		return domain.Order{}, fmt.Errorf("error while getting order with id %d: %w", id, err)
	}
	return order, nil
}
