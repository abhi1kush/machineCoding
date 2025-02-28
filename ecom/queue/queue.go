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
}

func NewQueue(poolSize int, queueCapacity int, processOrderFunc func(item Item), metricRepo repository.MetricRepositoryI, orderRepo repository.OrderRepositoryI, cache cache.CacheI) *Queue {
	return &Queue{
		orderQueue:       make(chan Item, queueCapacity),
		workerPool:       poolSize,
		wg:               sync.WaitGroup{},
		orderRepo:        orderRepo,
		metricRepo:       metricRepo,
		processOrderFunc: processOrderFunc,
		cache:            cache,
	}
}

func (q *Queue) StartOrderProcessor() {
	for i := 0; i < q.workerPool; i++ {
		q.wg.Add(1)
		go q.worker(q.processOrderFunc)
	}
}

func (q *Queue) worker(processOrderFunc func(item Item)) {
	defer q.wg.Done()
	for item := range q.orderQueue {
		start := time.Now()
		processOrderFunc(item)
		duration := time.Since(start)
		err := q.metricRepo.CreateMetric(
			&models.Metric{
				OrderId:    item.Id,
				Duration:   duration.Seconds(),
				MetricName: string(constants.PROCESSING_TIME),
			})
		if err != nil {
			log.Println("Error updating metrics in MetricsDB:", err)
		}
	}
}

func (q *Queue) AddOrderToQueue(item Item) {
	q.orderQueue <- item
}
