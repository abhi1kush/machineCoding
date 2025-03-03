package services

import (
	"database/sql"
	"log"
	"sync"
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
	itemRepo             repository.ItemRepositoryI
	orderCreationQueue   queue.QueueI
	orderProcessingQueue queue.QueueI
	cache                cache.CacheI
}

func NewOrderService(appConfig config.Config, orderRepo repository.OrderRepositoryI, itemRepo repository.ItemRepositoryI, metricRepo repository.MetricRepositoryI, cache cache.CacheI) *Order {
	orderService := &Order{
		repo:     orderRepo,
		itemRepo: itemRepo,
		cache:    cache,
	}
	orderService.orderCreationQueue = queue.NewQueue(appConfig.Queue.WorkerPool, appConfig.Queue.QueueCapacity, orderService.CreateOrderInDB, metricRepo, orderRepo, cache)
	orderService.orderProcessingQueue = queue.NewQueue(appConfig.Queue.WorkerPool, appConfig.Queue.QueueCapacity, orderService.ProcessOrder, metricRepo, orderRepo, cache)
	return orderService
}

const (
	ProcessingTimeMetricKey = "processing_time"
	CreationTimeMetricKey   = "creation_time"
)

func (o *Order) CreateOrder(userID string, itemIDs []string, totalAmount float64) (string, error) {
	orderID := uuid.New().String()

	err := o.cache.SetOrderStatus(orderID, string(constants.PENDING))
	if err != nil {
		log.Printf("Warning: failed to set order status in Redis for orderID %s: %v", orderID, err)
	}

	o.orderCreationQueue.Enqueue(queue.Item{Id: orderID, Value: &common.OrderRequest{UserID: userID, ItemIDs: itemIDs, TotalAmount: totalAmount}})
	return orderID, nil
}

func (o *Order) GetOrder(orderID string) (*common.OrderResponse, error) {
	return o.getOrder(orderID)
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
	// Simulating Order Process Delay.
	time.Sleep(1 * time.Second)
	//OrderProcess completed.
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(orderID string, wg *sync.WaitGroup) {
		defer wg.Done()
		if err := o.repo.UpdateOrderStatus(orderID, string(constants.COMPELETED)); err != nil {
			log.Println("Error updating order to Completed in DB:", err)
		}
	}(order.OrderID, &wg)

	if err := o.cache.SetOrderStatus(order.OrderID, string(constants.COMPELETED)); err != nil {
		log.Printf("Error updating cache ,order %v to Completed: err %v", order.OrderID, err)
	}
	wg.Wait()
}

func (o *Order) CreateOrderInDB(qItem queue.Item) {
	orderReq, _ := qItem.Value.(*common.OrderRequest)
	err := o.saveOrderInDB(qItem.Id, *orderReq)
	if err != nil {
		log.Printf("Failed to saveOrderInDb %v err %v", qItem.Id, err)
	}
	o.orderProcessingQueue.Enqueue(queue.Item{Id: qItem.Id, Value: &common.OrderItem{OrderID: qItem.Id}})
}

func (o *Order) GetOrderProcessQueue() queue.QueueI {
	return o.orderProcessingQueue
}

func (o *Order) GetOrderCreationQueue() queue.QueueI {
	return o.orderCreationQueue
}

func (o *Order) saveOrderInDB(orderId string, req common.OrderRequest) error {
	err := o.repo.CreateOrder(&models.Order{
		OrderID:     orderId,
		UserID:      req.UserID,
		TotalAmount: req.TotalAmount,
		Status:      string(constants.PENDING),
	})
	if err != nil {
		return err
	}
	for _, itemId := range req.ItemIDs {
		err := o.itemRepo.CreateItem(&models.Item{ItemID: itemId, OrderID: orderId})
		if err != nil {
			log.Printf("Failed to create Item %v err %v", itemId, err)
		}
	}
	return nil
}

func (o *Order) getOrder(orderID string) (*common.OrderResponse, error) {
	order, err := o.repo.GetOrderByID(orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	items, err := o.itemRepo.GetItemsByOrderId(orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	var itemIds []string
	for _, item := range items {
		itemIds = append(itemIds, item.ItemID)
	}
	orderResp := &common.OrderResponse{
		OrderID:     order.OrderID,
		UserID:      order.UserID,
		ItemIDs:     itemIds,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
	}
	return orderResp, nil
}
