package queue

import (
	"log"
	"sync"
	"time"

	"ecom.com/cache"
	"ecom.com/constants"
	"ecom.com/db"
)

type Order struct {
	OrderID string
}

var (
	orderQueue chan Order
	workerPool int
	wg         sync.WaitGroup
)

func StartOrderProcessor(poolSize int, queueCapacity int) {
	orderQueue = make(chan Order, queueCapacity)
	workerPool = poolSize

	for i := 0; i < workerPool; i++ {
		wg.Add(1)
		go worker()
	}
}

func worker() {
	defer wg.Done()
	for order := range orderQueue {
		start := time.Now()
		processOrder(order)
		duration := time.Since(start)

		_, err := db.MetricsDB.Exec("INSERT INTO metrics (order_id, processing_time) VALUES (?, ?)", order.OrderID, duration.Seconds())
		if err != nil {
			log.Println("Error updating metrics in MetricsDB:", err)
		}
	}
}

func AddOrderToQueue(orderID string) {
	orderQueue <- Order{OrderID: orderID}
}

func processOrder(order Order) {
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
