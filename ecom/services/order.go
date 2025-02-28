package services

import (
	"database/sql"
	"log"
	"time"

	"ecom.com/cache"
	"ecom.com/common"
	"ecom.com/constants"
	"ecom.com/db"
	"ecom.com/logger"
	"ecom.com/queue"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

const (
	ProcessingTimeMetricKey = "processing_time"
	CreationTimeMetricKey   = "creation_time"
)

func CreateOrder(userID string, itemIDs string, totalAmount float64) (string, error) {
	orderID := uuid.New().String()

	err := cache.SetOrderStatus(orderID, string(constants.PENDING))
	if err != nil {
		log.Printf("Warning: failed to set order status in Redis for orderID %s: %v", orderID, err)
	}

	queue.OrderCreationQueue.AddOrderToQueue(queue.Item{Id: orderID, Value: &common.OrderRequest{UserID: userID, ItemIDs: itemIDs, TotalAmount: totalAmount}})
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

func ProcessOrder(item queue.Item) {
	order, ok := item.Value.(*common.OrderItem)
	if !ok {
		log.Println("Invalid item in queue:")
		return
	}
	if err := cache.SetOrderStatus(order.OrderID, string(constants.PROCESSING)); err != nil {
		log.Println("Error updating cache to Processing:", err)
	}

	time.Sleep(1 * time.Second)

	if err := cache.SetOrderStatus(order.OrderID, string(constants.COMPELETED)); err != nil {
		log.Printf("Error updating cache ,order %v to Completed: err %v", order.OrderID, err)
	}

	go func(orderID string) {
		if _, err := db.DB.Exec("UPDATE orders SET status = ? WHERE order_id = ?", string(constants.COMPELETED), orderID); err != nil {
			log.Println("Error updating order to Completed in DB:", err)
		}
	}(order.OrderID)
}

func CreateOrderInDB(item queue.Item) {
	orderReq, _ := item.Value.(*common.OrderRequest)
	_, err := db.DB.Exec("INSERT INTO orders (order_id, user_id, item_ids, total_amount, status) VALUES (?, ?, ?, ?, ?)",
		item.Id, orderReq.UserID, orderReq.ItemIDs, orderReq.TotalAmount, constants.PENDING)
	if err != nil {
		log.Printf("Failed to create order %v", err)
	}
	queue.OrderProcessingQueue.AddOrderToQueue(queue.Item{Id: item.Id, Value: common.OrderItem{OrderID: item.Id}})
}
