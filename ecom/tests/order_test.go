package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

	"ecom.com/cache"
	"ecom.com/config"
	"ecom.com/db"
	"ecom.com/queue"
	"ecom.com/routes"
	"ecom.com/services"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	config.AppConfig.Server.Port = "8081"
	config.AppConfig.Database.Driver = "sqlite3"
	config.AppConfig.Database.DSN = "orders_test.db"
	config.AppConfig.Queue.WorkerPool = 5
	config.AppConfig.Queue.QueueCapacity = 500

	cache.InitRedis("127.0.0.1", "pass", 1)

	db.InitDB(config.AppConfig.Database.Driver, config.AppConfig.Database.DSN)
	defer db.CloseDB()

	db.InitMetricsDB(config.AppConfig.Database.Driver, config.AppConfig.Database.DSN)
	defer db.CloseDB()

	queue.OrderCreationQueue = queue.NewQueue(config.AppConfig.Queue.WorkerPool, config.AppConfig.Queue.QueueCapacity, services.CreateOrderInDB, services.CreationTimeMetricKey)
	queue.OrderCreationQueue.StartOrderProcessor()

	queue.OrderProcessingQueue = queue.NewQueue(config.AppConfig.Queue.WorkerPool, config.AppConfig.Queue.QueueCapacity, services.ProcessOrder, services.ProcessingTimeMetricKey)
	queue.OrderProcessingQueue.StartOrderProcessor()

	router = gin.Default()
	routes.SetupRoutes(router)

	m.Run()
}

func TestCreateOrder(t *testing.T) {
	orderPayload := map[string]interface{}{
		"user_id":      "test-user-123",
		"item_ids":     "item1,item2",
		"total_amount": 100.50,
	}
	payload, _ := json.Marshal(orderPayload)
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
	assert.Nil(t, err)
	assert.Contains(t, responseBody, "order_id")
}

func TestGetOrderStatusNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/orders/non-existent-id", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestCreateOrderInvalidPayload(t *testing.T) {
	orderPayload := map[string]interface{}{
		"user_id": 123,
	}
	payload, _ := json.Marshal(orderPayload)
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestLoad500ConcurrentRequests(t *testing.T) {
	var wg sync.WaitGroup
	requestCount := 500
	wg.Add(requestCount)
	start := time.Now()

	for i := 0; i < requestCount; i++ {
		go func(i int) {
			defer wg.Done()
			orderPayload := map[string]interface{}{
				"user_id":      "test-user-" + strconv.Itoa(i),
				"item_ids":     "item1,item2",
				"total_amount": 100.50,
			}
			payload, _ := json.Marshal(orderPayload)
			req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			assert.Equal(t, http.StatusOK, resp.Code)
		}(i)
	}
	wg.Wait()
	duration := time.Since(start)
	t.Logf("Completed %v concurrent requests in %s", requestCount, duration)
}

func TestGetMetrics(t *testing.T) {
	time.Sleep(3 * time.Second)
	req, _ := http.NewRequest("GET", "/metrics", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var metrics map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &metrics)
	assert.Nil(t, err)

	assert.Contains(t, metrics, "total_orders_received")
	assert.Contains(t, metrics, "average_processing_time")
	assert.Contains(t, metrics, "orders_pending")
	assert.Contains(t, metrics, "orders_processing")
	assert.Contains(t, metrics, "orders_completed")
}
