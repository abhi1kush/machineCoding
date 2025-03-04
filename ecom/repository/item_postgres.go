package repository

import (
	"database/sql"
	"log"

	"ecom.com/models"
)

type PostgreSqlItemRepository struct {
	DB *sql.DB
}

func NewPostgreSqlItemRepository(db *sql.DB) ItemRepositoryI {
	return &PostgreSqlItemRepository{DB: db}
}

func (r *PostgreSqlItemRepository) CreateItem(item *models.Item) error {
	itemQuery := `INSERT INTO items (item_id, order_id, amount) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(itemQuery, item.ItemID, item.OrderID, item.Amount)
	if err != nil {
		log.Printf("Failed to add item err %v", err)
	}
	return err
}

func (r *PostgreSqlItemRepository) GetItem(id string) (*models.Item, error) {
	query := `SELECT item_id, order_id, amount FROM items WHERE item_id = $1`
	row := r.DB.QueryRow(query, id)

	var item models.Item
	err := row.Scan(&item.ItemID, &item.OrderID, &item.Amount)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *PostgreSqlItemRepository) GetItemsByOrderId(id string) ([]models.Item, error) {
	query := `SELECT item_id, order_id, amount FROM items WHERE order_id = $1`
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

func (r *PostgreSqlItemRepository) RemoveItem(itemId string, orderId string) error {
	itemQuery := `DELETE FROM items WHERE item_id = $1 AND order_id = $2`
	_, err := r.DB.Exec(itemQuery, itemId, orderId)
	if err != nil {
		log.Printf("Failed to add item err %v", err)
	}
	return err
}
