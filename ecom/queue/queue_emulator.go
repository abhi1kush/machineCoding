package queue

import (
	"log"
	"sync"
	"time"

	"ecom.com/cache"
	"ecom.com/constants"
	"ecom.com/models"
	"ecom.com/repository"
)

type Item struct {
	Id    string
	Value any
}

type Queue struct {
	orderQueue       chan Item
	workerPool       int
	wg               sync.WaitGroup
	orderRepo        repository.OrderRepositoryI
	metricRepo       repository.MetricRepositoryI
	processOrderFunc func(item Item)
	cache            cache.CacheI
	stopChan         chan struct{} // Channel for graceful shutdown
}

func NewQueue(poolSize int, queueCapacity int, processOrderFunc func(item Item), metricRepo repository.MetricRepositoryI, orderRepo repository.OrderRepositoryI, cache cache.CacheI) QueueI {
	return &Queue{
		orderQueue:       make(chan Item, queueCapacity),
		workerPool:       poolSize,
		wg:               sync.WaitGroup{},
		orderRepo:        orderRepo,
		metricRepo:       metricRepo,
		processOrderFunc: processOrderFunc,
		cache:            cache,
		stopChan:         make(chan struct{}),
	}
}

func (q *Queue) StartOrderProcessor() error {
	for i := 0; i < q.workerPool; i++ {
		q.wg.Add(1)
		go q.worker(q.processOrderFunc)
	}
	return nil
}

func (q *Queue) worker(processOrderFunc func(item Item)) {
	defer q.wg.Done()
	for {
		select {
		case item, ok := <-q.orderQueue:
			if !ok {
				// Queue closed, exit worker.
				return
			}
			start := time.Now()
			processOrderFunc(item)
			duration := time.Since(start)

			// Log processing time as a metric
			err := q.metricRepo.CreateMetric(&models.Metric{
				OrderId:    item.Id,
				Duration:   duration.Seconds(),
				MetricName: string(constants.PROCESSING_TIME),
			})
			if err != nil {
				log.Println("Error updating metrics in MetricsDB:", err)
			}
		case <-q.stopChan:
			// Received stop signal, exit gracefully.
			return
		}
	}
}

func (q *Queue) Enqueue(item Item) {
	select {
	case q.orderQueue <- item:
		// Successfully enqueued item
	default:
		log.Println("Warning: Queue is full. Dropping item:", item.Id)
	}
}

// StopOrderProcessor gracefully shuts down all workers.
func (q *Queue) StopOrderProcessor() {
	close(q.stopChan)   // Notify workers to stop
	close(q.orderQueue) // Close queue to prevent new items

	q.wg.Wait() // Wait for all workers to finish
	log.Println("Order processing stopped.")
}
