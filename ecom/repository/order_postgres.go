package repository

import (
	"database/sql"

	"ecom.com/models"
)

type PostgreSqlOrderRepository struct {
	DB *sql.DB
}

func NewPostgreSqlOrderRepository(db *sql.DB) OrderRepositoryI {
	return &PostgreSqlOrderRepository{DB: db}
}

func (r *PostgreSqlOrderRepository) CreateOrder(order *models.Order) error {
	query := `INSERT INTO orders (order_id, user_id, total_amount, status) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(query, order.OrderID, order.UserID, order.TotalAmount, order.Status)
	return err
}

func (r *PostgreSqlOrderRepository) UpdateOrderStatus(orderId string, status string) error {
	query := `UPDATE orders SET status = $1 WHERE order_id = $2;`
	_, err := r.DB.Exec(query, status, orderId)
	return err
}

func (r *PostgreSqlOrderRepository) GetOrderByID(id string) (*models.Order, error) {
	query := `SELECT order_id, user_id, total_amount, status FROM orders WHERE order_id = $1`
	row := r.DB.QueryRow(query, id)

	var order models.Order
	err := row.Scan(&order.OrderID, &order.UserID, &order.TotalAmount, &order.Status)
	if err != nil {
		return nil, err
	}
	return &order, nil
}
