package services

import (
	"database/sql"
	"log"

	"ecom.com/cache"
	"ecom.com/constants"
	"ecom.com/db"
	"ecom.com/logger"
	"ecom.com/queue"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

func CreateOrder(userID string, itemIDs string, totalAmount float64) (string, error) {
	orderID := uuid.New().String()

	_, err := db.DB.Exec("INSERT INTO orders (order_id, user_id, item_ids, total_amount, status) VALUES (?, ?, ?, ?, ?)",
		orderID, userID, itemIDs, totalAmount, constants.PENDING)
	if err != nil {
		return "", err
	}

	err = cache.SetOrderStatus(orderID, string(constants.PENDING))
	if err != nil {
		log.Printf("Warning: failed to set order status in Redis for orderID %s: %v", orderID, err)
	}

	queue.AddOrderToQueue(orderID)
	return orderID, nil
}

func GetOrderStatus(orderID string) (string, error) {
	status, err := cache.GetOrderStatus(orderID)
	if err == nil {
		return status, nil
	}

	if err != redis.Nil {
		logger.Logger.Printf("Redis error %v", err)
	}

	// Fallback to DB.
	var dbStatus string
	err = db.DB.QueryRow("SELECT status FROM orders WHERE order_id = ?", orderID).Scan(&dbStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", sql.ErrNoRows
		}
		return "", err
	}

	_ = cache.SetOrderStatus(orderID, dbStatus)

	return dbStatus, nil
}
