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
	"ecom.com/queue"

	"github.com/stretchr/testify/assert"
)

// TestCreateOrderAPI tests that a valid order is created via the API and that an order_id is returned.
func TestCreateOrderAPI(t *testing.T) {
	orderPayload := map[string]interface{}{
		"user_id":      "test-user-api",
		"item_ids":     []string{"item1", "item2"},
		"total_amount": 123.45,
	}
	payloadBytes, _ := json.Marshal(orderPayload)
	req, _ := http.NewRequest("POST", "/api/v1/orders", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	globalTestRouter.ServeHTTP(w, req)

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
	// Create order using the service, which also sets Redis.
	orderID, err := globalTestContainer.OrderService.CreateOrder("test-user-status", []string{"item1", "item2"}, 50.0)
	assert.Nil(t, err)

	var response map[string]interface{}
	// Check initial status via API.
	req, _ := http.NewRequest("GET", "/api/v1/orders/status/"+orderID, nil)
	w := httptest.NewRecorder()
	globalTestRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Pending", response["status"])

	time.Sleep(3 * time.Second)

	// Retrieve status again; final status should be "Completed".
	req2, _ := http.NewRequest("GET", "/api/v1/orders/status/"+orderID, nil)
	w2 := httptest.NewRecorder()
	globalTestRouter.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	err = json.Unmarshal(w2.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Completed", response["status"])
}

// TestMetricsAPI tests the /metrics endpoint to ensure it returns valid metrics.
func TestMetricsAPI(t *testing.T) {
	// Create and process a few orders.
	count := 3
	for i := 0; i < count; i++ {
		_, err := globalTestContainer.OrderService.CreateOrder("user"+strconv.Itoa(i), []string{"item1", "item2"}, 100.0)
		assert.Nil(t, err)
	}
	// Allow processing to complete.
	time.Sleep(4 * time.Second)

	req, _ := http.NewRequest("GET", "/api/v1/metrics", nil)
	w := httptest.NewRecorder()
	globalTestRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	totalOrdersProcessed := int(response["total_orders_received"].(float64))
	assert.GreaterOrEqual(t, totalOrdersProcessed, count)
}

func TestDatabaseOperations0(t *testing.T) {
	orderID := 5
	userID := "db-test-user"
	// itemIDs := []string{"item1", "item2"}
	amount := 75.0
	var status string
	_, err := globalTestContainer.DB.Exec("INSERT INTO orders (order_id, user_id, total_amount, status) VALUES (?, ?, ?, ?)", orderID, userID, amount, "Pending")
	assert.Nil(t, err)
	var orderIDOut, userId string
	var amountOut float64
	err = globalTestContainer.DB.QueryRow("SELECT order_id, user_id, total_amount, status  FROM orders WHERE order_id = ?", orderID).Scan(&orderIDOut, &userId, &amountOut, &status)
	assert.Nil(t, err)
	assert.Equal(t, "Pending", status)
	assert.Equal(t, amount, amountOut)
	assert.Equal(t, "db-test-user", userId)
}

// TestDatabaseOperations verifies that creating an order inserts the correct record into the orders DB.
func TestDatabaseOperations(t *testing.T) {
	orderID, err := globalTestContainer.OrderService.CreateOrder("db-test-user", []string{"item1", "item2"}, 75.0)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)
	var status string
	var orderIDOut, userId string
	var amount float64
	err = globalTestContainer.DB.QueryRow("SELECT order_id, user_id, total_amount, status  FROM orders WHERE order_id = ?", orderID).Scan(&orderIDOut, &userId, &amount, &status)
	assert.Nil(t, err)
	assert.Equal(t, 75.0, amount)
	assert.Equal(t, "db-test-user", userId)
	assert.Equal(t, "Pending", status)
}

// TestQueueProcessing tests that an order added to the queue is processed and updated to "Completed"
// in both the orders DB and the Redis cache.
func TestQueueProcessing(t *testing.T) {
	orderID, err := globalTestContainer.OrderService.CreateOrder("queue-test-user", []string{"item1", "item2"}, 200.0)
	assert.Nil(t, err)
	time.Sleep(4 * time.Second)

	// Verify orders DB status.
	var status string
	err = globalTestContainer.DB.QueryRow("SELECT status FROM orders WHERE order_id = ?", orderID).Scan(&status)
	assert.Nil(t, err)
	assert.Equal(t, "Completed", status)

	// Verify Redis cache status.
	cachedStatus, err := globalTestContainer.Cache.GetOrderStatus(orderID)
	assert.Nil(t, err)
	assert.Equal(t, "Completed", cachedStatus)
}

// TestConcurrentQueueProcessing tests the queue's ability to process multiple orders concurrently.
func TestConcurrentQueueProcessing(t *testing.T) {
	var wg sync.WaitGroup
	numOrders := 15
	for i := 0; i < numOrders; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			orderID, err := globalTestContainer.OrderService.CreateOrder("concurrent-user"+strconv.Itoa(i), []string{"item1", "item2"}, 150.0)
			if err != nil {
				t.Errorf("Error creating order: %v", err)
				return
			}
			globalTestContainer.OrderService.GetOrderProcessQueue().Enqueue(queue.Item{Id: orderID, Value: &common.OrderItem{OrderID: orderID}})
		}(i)
	}
	wg.Wait()

	// Wait for all orders to be processed.
	time.Sleep(10 * time.Second)

	// Verify each order is marked as "Completed" in the orders DB.
	rows, err := globalTestContainer.DB.Query("SELECT status FROM orders")
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var status string
		err = rows.Scan(&status)
		assert.Nil(t, err)
		assert.Equal(t, "Completed", status)
	}
}
