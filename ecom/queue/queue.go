package queue

import (
	"log"
	"sync"
	"time"

	"ecom.com/db"
)

var OrderProcessingQueue *Queue
var OrderCreationQueue *Queue

type Item struct {
	Id    string
	Value any
}

type Queue struct {
	orderQueue       chan Item
	workerPool       int
	wg               sync.WaitGroup
	processOrderFunc func(item Item)
	metricKeyName    string
}

func NewQueue(poolSize int, queueCapacity int, processOrderFunc func(item Item), metricKeyName string) *Queue {
	return &Queue{
		orderQueue: make(chan Item, queueCapacity),
		workerPool: poolSize, wg: sync.WaitGroup{},
		processOrderFunc: processOrderFunc,
		metricKeyName:    metricKeyName,
	}
}

func (q *Queue) StartOrderProcessor() {
	for i := 0; i < q.workerPool; i++ {
		q.wg.Add(1)
		go q.worker(q.processOrderFunc, q.metricKeyName)
	}
}

func (q *Queue) worker(processOrderFunc func(item Item), metricKeyName string) {
	defer q.wg.Done()
	for item := range q.orderQueue {
		start := time.Now()
		processOrderFunc(item)
		duration := time.Since(start)
		_, err := db.MetricsDB.Exec("INSERT INTO metrics (order_id, processing_time) VALUES (?, ?)", item.Id, duration.Seconds())
		if err != nil {
			log.Println("Error updating metrics in MetricsDB:", err)
		}
	}
}

func (q *Queue) AddOrderToQueue(item Item) {
	q.orderQueue <- item
}
