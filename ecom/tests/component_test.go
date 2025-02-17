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

	"ecom.com/config"
	"ecom.com/db"
	"ecom.com/queue"
	"ecom.com/routes"
	"ecom.com/services"

	"ecom.com/cache"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var testRouter *gin.Engine

func setup() {
	config.AppConfig.Server.Port = "8081"
	config.AppConfig.Database.Driver = "sqlite3"
	config.AppConfig.Database.DSN = "orders_test.db"
	config.AppConfig.Metrics.Driver = "sqlite3"
	config.AppConfig.Metrics.DSN = "metrics_test.db"
	config.AppConfig.Queue.WorkerPool = 5
	config.AppConfig.Queue.QueueCapacity = 500
	config.AppConfig.Redis.Addr = "localhost:6379"
	config.AppConfig.Redis.Password = ""
	config.AppConfig.Redis.DB = 1

	gin.SetMode(gin.TestMode)

	db.InitDB(config.AppConfig.Database.Driver, config.AppConfig.Database.DSN)
	db.InitMetricsDB(config.AppConfig.Metrics.Driver, config.AppConfig.Metrics.DSN)
	db.DB.Exec("DELETE FROM orders")
	db.MetricsDB.Exec("DELETE FROM metrics")

	cache.InitRedis(config.AppConfig.Redis.Addr, config.AppConfig.Redis.Password, config.AppConfig.Redis.DB)

	queue.StartOrderProcessor(config.AppConfig.Queue.WorkerPool, config.AppConfig.Queue.QueueCapacity)

	testRouter = gin.Default()
	routes.SetupRoutes(testRouter)
}

func teardown() {
	db.CloseDB()
	db.CloseMetricsDB()
}

// TestCreateOrderAPI tests that a valid order is created via the API and that an order_id is returned.
func TestCreateOrderAPI(t *testing.T) {
	setup()
	defer teardown()

	orderPayload := map[string]interface{}{
		"user_id":      "test-user-api",
		"item_ids":     "item1,item2",
		"total_amount": 123.45,
	}
	payloadBytes, _ := json.Marshal(orderPayload)
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	_, ok := response["order_id"]
	assert.True(t, ok, "order_id should be returned")
}

// TestGetOrderStatusAPI tests that an order's status is retrieved correctly.
// It verifies that the status is initially "Pending" then transitions to "Completed" after processing.
func TestGetOrderStatusAPI(t *testing.T) {
	setup()
	defer teardown()

	// Create order using the service, which also sets Redis.
	orderID, err := services.CreateOrder("test-user-status", "item1,item2", 50.0)
	assert.Nil(t, err)

	// Check initial status via API.
	req, _ := http.NewRequest("GET", "/order/"+orderID, nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Pending", response["status"])

	// Add order to processing queue and wait for processing.
	queue.AddOrderToQueue(orderID)
	time.Sleep(3 * time.Second)

	// Retrieve status again; final status should be "Completed".
	req2, _ := http.NewRequest("GET", "/order/"+orderID, nil)
	w2 := httptest.NewRecorder()
	testRouter.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	err = json.Unmarshal(w2.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Completed", response["status"])
}

// TestMetricsAPI tests the /metrics endpoint to ensure it returns valid metrics.
func TestMetricsAPI(t *testing.T) {
	setup()
	defer teardown()

	// Create and process a few orders.
	count := 3
	for i := 0; i < count; i++ {
		orderID, err := services.CreateOrder("user"+strconv.Itoa(i), "item1,item2", 100.0)
		assert.Nil(t, err)
		queue.AddOrderToQueue(orderID)
	}
	// Allow processing to complete.
	time.Sleep(4 * time.Second)

	req, _ := http.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	totalOrdersProcessed := int(response["total_orders_received"].(float64))
	assert.GreaterOrEqual(t, totalOrdersProcessed, count)
}

// TestDatabaseOperations verifies that creating an order inserts the correct record into the orders DB.
func TestDatabaseOperations(t *testing.T) {
	setup()
	defer teardown()

	orderID, err := services.CreateOrder("db-test-user", "item1,item2", 75.0)
	assert.Nil(t, err)

	var status string
	err = db.DB.QueryRow("SELECT status FROM orders WHERE order_id = ?", orderID).Scan(&status)
	assert.Nil(t, err)
	assert.Equal(t, "Pending", status)
}

// TestQueueProcessing tests that an order added to the queue is processed and updated to "Completed"
// in both the orders DB and the Redis cache.
func TestQueueProcessing(t *testing.T) {
	setup()
	defer teardown()

	orderID, err := services.CreateOrder("queue-test-user", "item1,item2", 200.0)
	assert.Nil(t, err)
	queue.AddOrderToQueue(orderID)
	time.Sleep(4 * time.Second)

	// Verify orders DB status.
	var status string
	err = db.DB.QueryRow("SELECT status FROM orders WHERE order_id = ?", orderID).Scan(&status)
	assert.Nil(t, err)
	assert.Equal(t, "Completed", status)

	// Verify Redis cache status.
	cachedStatus, err := cache.GetOrderStatus(orderID)
	assert.Nil(t, err)
	assert.Equal(t, "Completed", cachedStatus)
}

// TestConcurrentQueueProcessing tests the queue's ability to process multiple orders concurrently.
func TestConcurrentQueueProcessing(t *testing.T) {
	setup()
	defer teardown()

	var wg sync.WaitGroup
	numOrders := 20
	for i := 0; i < numOrders; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			orderID, err := services.CreateOrder("concurrent-user"+strconv.Itoa(i), "item1,item2", 150.0)
			if err != nil {
				t.Errorf("Error creating order: %v", err)
				return
			}
			queue.AddOrderToQueue(orderID)
		}(i)
	}
	wg.Wait()

	// Wait for all orders to be processed.
	time.Sleep(10 * time.Second)

	// Verify each order is marked as "Completed" in the orders DB.
	rows, err := db.DB.Query("SELECT status FROM orders")
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var status string
		err = rows.Scan(&status)
		assert.Nil(t, err)
		assert.Equal(t, "Completed", status)
	}
}
