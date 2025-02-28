package repository

import (
	"database/sql"

	"ecom.com/models"
)

type SQLiteOrderRepository struct {
	DB *sql.DB
}

func NewSQLiteOrderRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{DB: db}
}

func (r *SQLiteUserRepository) CreateOrder(order *models.Order) error {
	query := `INSERT INTO orders (order_id, user_id, item_ids, total_amount, status) VALUES (?, ?, ?, ?, ?)`
	_, err := r.DB.Exec(query, order.OrderID, order.UserID, order.ItemIDs, order.TotalAmount, order.Status)
	return err
}

func (r *SQLiteUserRepository) UpdateOrderStatus(orderId string, status string) error {
	query := `UPDATE a SET a.status = ? FROM orders WHERE a.order_id = ?`
	_, err := r.DB.Exec(query, orderId, status)
	return err
}

func (r *SQLiteUserRepository) GetOrderByID(id string) (*models.Order, error) {
	query := `SELECT order_id, user_id, total_amount, status FROM orders WHERE order_id = ?`
	row := r.DB.QueryRow(query, id)

	var order models.Order
	err := row.Scan(&order.OrderID, &order.UserID, &order.TotalAmount, &order.Status)
	if err != nil {
		return nil, err
	}
	return &order, nil
}
