package services

import (
	"database/sql"
	"log"
	"time"

	"ecom.com/cache"
	"ecom.com/common"
	"ecom.com/config"
	"ecom.com/constants"
	"ecom.com/logger"
	"ecom.com/models"
	"ecom.com/queue"
	"ecom.com/repository"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

type Order struct {
	repo                 repository.OrderRepositoryI
	orderCreationQueue   *queue.Queue
	orderProcessingQueue *queue.Queue
	cache                cache.CacheI
}

func NewOrderService(appConfig config.Config, orderRepo repository.OrderRepositoryI, metricRepo repository.MetricRepositoryI, cache cache.CacheI) *Order {
	orderService := &Order{
		repo:  orderRepo,
		cache: cache,
	}
	orderService.orderCreationQueue = queue.NewQueue(appConfig.Queue.WorkerPool, appConfig.Queue.QueueCapacity, orderService.CreateOrderInDB, metricRepo, orderRepo, cache)
	orderService.orderProcessingQueue = queue.NewQueue(appConfig.Queue.WorkerPool, appConfig.Queue.QueueCapacity, orderService.ProcessOrder, metricRepo, orderRepo, cache)
	return orderService
}

const (
	ProcessingTimeMetricKey = "processing_time"
	CreationTimeMetricKey   = "creation_time"
)

func (o *Order) CreateOrder(userID string, itemIDs string, totalAmount float64) (string, error) {
	orderID := uuid.New().String()

	err := o.cache.SetOrderStatus(orderID, string(constants.PENDING))
	if err != nil {
		log.Printf("Warning: failed to set order status in Redis for orderID %s: %v", orderID, err)
	}

	o.orderCreationQueue.AddOrderToQueue(queue.Item{Id: orderID, Value: &common.OrderRequest{UserID: userID, ItemIDs: itemIDs, TotalAmount: totalAmount}})
	return orderID, nil
}

func (o *Order) GetOrder(orderID string) (*common.OrderResponse, error) {
	order, err := o.repo.GetOrderByID(orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	orderResp := &common.OrderResponse{
		OrderID:     order.OrderID,
		UserID:      order.UserID,
		ItemIDs:     order.ItemIDs,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
	}
	return orderResp, nil
}

func (o *Order) GetOrderStatus(orderID string) (string, error) {
	status, err := o.cache.GetOrderStatus(orderID)
	if err == nil {
		return status, nil
	}

	if err != redis.Nil {
		logger.Logger.Printf("Redis error %v", err)
	}

	// Fallback to DB.
	order, err := o.repo.GetOrderByID(orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", sql.ErrNoRows
		}
		return "", err
	}

	_ = o.cache.SetOrderStatus(order.OrderID, order.Status)
	return order.Status, nil
}

func (o *Order) ProcessOrder(item queue.Item) {
	order, ok := item.Value.(*common.OrderItem)
	if !ok {
		log.Printf("Invalid item in queue: %v ", item)
		return
	}
	if err := o.cache.SetOrderStatus(order.OrderID, string(constants.PROCESSING)); err != nil {
		log.Println("Error updating cache to Processing:", err)
	}

	time.Sleep(1 * time.Second)

	if err := o.cache.SetOrderStatus(order.OrderID, string(constants.COMPELETED)); err != nil {
		log.Printf("Error updating cache ,order %v to Completed: err %v", order.OrderID, err)
	}

	go func(orderID string) {
		if err := o.repo.UpdateOrderStatus(orderID, string(constants.COMPELETED)); err != nil {
			log.Println("Error updating order to Completed in DB:", err)
		}
	}(order.OrderID)
}

func (o *Order) CreateOrderInDB(item queue.Item) {
	orderReq, _ := item.Value.(*common.OrderRequest)
	err := o.repo.CreateOrder(&models.Order{
		OrderID:     item.Id,
		UserID:      orderReq.UserID,
		ItemIDs:     orderReq.ItemIDs,
		TotalAmount: orderReq.TotalAmount,
		Status:      string(constants.PENDING),
	})
	if err != nil {
		log.Printf("Failed to create order %v", err)
	}
	o.orderProcessingQueue.AddOrderToQueue(queue.Item{Id: item.Id, Value: common.OrderItem{OrderID: item.Id}})
}

func (o *Order) GetOrderProcessQueue() *queue.Queue {
	return o.orderProcessingQueue
}

func (o *Order) GetOrderCreationQueue() *queue.Queue {
	return o.orderCreationQueue
}
