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

	"ecom.com/common"
	"ecom.com/config"
	"ecom.com/database"
	"ecom.com/queue"
	"ecom.com/routes"
	"ecom.com/server"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var testRouter *gin.Engine
var testContainer *server.Container

func setup() {
	testConfig := config.Config{}
	testConfig.Server.Port = "8081"
	testConfig.Database.Driver = "sqlite3"
	testConfig.Database.DSN = "orders_test.db"
	testConfig.Metrics.Driver = "sqlite3"
	testConfig.Metrics.DSN = "metrics_test.db"
	testConfig.Queue.WorkerPool = 5
	testConfig.Queue.QueueCapacity = 500
	testConfig.Redis.Addr = "localhost:6379"
	testConfig.Redis.Password = ""
	testConfig.Redis.DB = 1

	gin.SetMode(gin.TestMode)

	testContainer = server.NewContainer(testConfig)

	testContainer.OrderService.GetOrderCreationQueue().StartOrderProcessor()
	testContainer.OrderService.GetOrderProcessQueue().StartOrderProcessor()

	r := gin.Default()
	routes.RegisterRoutes(r, testContainer.RoutesCfg)
	r.Run()
}

func teardown() {
	database.CloseDB(testContainer.DB)
	database.CloseDB(testContainer.MetricDB)
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
	orderID, err := testContainer.OrderService.CreateOrder("test-user-status", "item1,item2", 50.0)
	assert.Nil(t, err)

	var response map[string]interface{}
	// Check initial status via API.
	req, _ := http.NewRequest("GET", "/order/"+orderID, nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Pending", response["status"])

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
		_, err := testContainer.OrderService.CreateOrder("user"+strconv.Itoa(i), "item1,item2", 100.0)
		assert.Nil(t, err)
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

func TestDatabaseOperations0(t *testing.T) {
	setup()
	defer teardown()
	orderID := 5
	userID := "db-test-user"
	itemIDs := "item1,item2"
	amount := 75.0
	var status string
	_, err := testContainer.DB.Exec("INSERT INTO orders (order_id, user_id, item_ids, total_amount, status) VALUES (?, ?, ?, ?, ?)", orderID, userID, itemIDs, amount, "Pending")
	assert.Nil(t, err)
	var orderIDOut, userId string
	var amountOut float64
	err = testContainer.DB.QueryRow("SELECT order_id, user_id, total_amount, status  FROM orders WHERE order_id = ?", orderID).Scan(&orderIDOut, &userId, &amountOut, &status)
	assert.Nil(t, err)
	assert.Equal(t, "Pending", status)
	assert.Equal(t, amount, amountOut)
	assert.Equal(t, "db-test-user", userId)
}

// TestDatabaseOperations verifies that creating an order inserts the correct record into the orders DB.
func TestDatabaseOperations(t *testing.T) {
	setup()
	defer teardown()

	orderID, err := testContainer.OrderService.CreateOrder("db-test-user", "item1,item2", 75.0)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)
	var status string
	var orderIDOut, userId string
	var amount float64
	err = testContainer.DB.QueryRow("SELECT order_id, user_id, total_amount, status  FROM orders WHERE order_id = ?", orderID).Scan(&orderIDOut, &userId, &amount, &status)
	assert.Nil(t, err)
	assert.Equal(t, 75.0, amount)
	assert.Equal(t, "db-test-user", userId)
	assert.Equal(t, "Pending", status)
}

// TestQueueProcessing tests that an order added to the queue is processed and updated to "Completed"
// in both the orders DB and the Redis cache.
func TestQueueProcessing(t *testing.T) {
	setup()
	defer teardown()

	orderID, err := testContainer.OrderService.CreateOrder("queue-test-user", "item1,item2", 200.0)
	assert.Nil(t, err)
	time.Sleep(4 * time.Second)

	// Verify orders DB status.
	var status string
	err = testContainer.DB.QueryRow("SELECT status FROM orders WHERE order_id = ?", orderID).Scan(&status)
	assert.Nil(t, err)
	assert.Equal(t, "Completed", status)

	// Verify Redis cache status.
	cachedStatus, err := testContainer.Cache.GetOrderStatus(orderID)
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
			orderID, err := testContainer.OrderService.CreateOrder("concurrent-user"+strconv.Itoa(i), "item1,item2", 150.0)
			if err != nil {
				t.Errorf("Error creating order: %v", err)
				return
			}
			testContainer.OrderService.GetOrderProcessQueue().AddOrderToQueue(queue.Item{Id: orderID, Value: common.OrderItem{OrderID: orderID}})
		}(i)
	}
	wg.Wait()

	// Wait for all orders to be processed.
	time.Sleep(10 * time.Second)

	// Verify each order is marked as "Completed" in the orders DB.
	rows, err := testContainer.DB.Query("SELECT status FROM orders")
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var status string
		err = rows.Scan(&status)
		assert.Nil(t, err)
		assert.Equal(t, "Completed", status)
	}
}
