package repository

import (
	"database/sql"
	"log"

	"ecom.com/models"
)

type SQLiteItemRepository struct {
	DB *sql.DB
}

func NewSQLiteItemRepository(db *sql.DB) ItemRepositoryI {
	return &SQLiteItemRepository{DB: db}
}

func (r *SQLiteItemRepository) CreateItem(item *models.Item) error {
	itemQuery := `INSERT INTO items (item_id, order_id, amount) VALUES (?, ?, ?)`
	_, err := r.DB.Exec(itemQuery, item.ItemID, item.OrderID, item.Amount)
	if err != nil {
		log.Printf("Failed to add item err %v", err)
	}
	return err
}

func (r *SQLiteItemRepository) GetItem(id string) (*models.Item, error) {
	query := `SELECT item_id, order_id, amount FROM items WHERE item_id = ?`
	row := r.DB.QueryRow(query, id)

	var item models.Item
	err := row.Scan(&item.ItemID, &item.OrderID, &item.Amount)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *SQLiteItemRepository) GetItemsByOrderId(id string) ([]models.Item, error) {
	query := `SELECT item_id, order_id, amount FROM items WHERE order_id = ?`
	rows, err := r.DB.Query(query, id)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()
	var items []models.Item
	for rows.Next() {
		item := models.Item{}
		err := rows.Scan(&item.ItemID, &item.OrderID, &item.Amount)
		if err != nil {
			log.Printf("GetItemsByOrderId: Failed parse this row err %v", err)
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *SQLiteItemRepository) RemoveItem(itemId string, orderId string) error {
	itemQuery := `DELETE FROM items WHERE item_id = ? AND order_id = ?`
	_, err := r.DB.Exec(itemQuery, itemId, orderId)
	if err != nil {
		log.Printf("Failed to add item err %v", err)
	}
	return err
}
