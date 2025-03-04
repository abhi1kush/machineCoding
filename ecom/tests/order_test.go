package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"ecom.com/config"
	"ecom.com/database"
	"ecom.com/routes"
	"ecom.com/server"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var globalTestRouter *gin.Engine
var globalTestContainer *server.Container

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	testConfig := config.Config{}
	testConfig.Server.Port = "8081"
	testConfig.Database.Driver = "sqlite3"
	testConfig.Metrics.Driver = "sqlite3"
	testConfig.Database.DSN = "orders_test.db"
	testConfig.Metrics.DSN = "metrics_test.db"
	testConfig.Queue.WorkerPool = 5
	testConfig.Queue.QueueCapacity = 500
	testConfig.Redis.Addr = "localhost:6379"
	testConfig.Redis.Password = ""
	testConfig.Redis.DB = 1

	globalTestContainer = server.NewContainer(testConfig)
	defer database.CloseDB(globalTestContainer.DB)
	defer database.CloseDB(globalTestContainer.MetricDB)

	globalTestContainer.OrderService.GetOrderProcessQueue().StartOrderProcessor()
	globalTestContainer.OrderService.GetOrderCreationQueue().StartOrderProcessor()

	// Initialize and assign router
	globalTestRouter = gin.Default()
	routes.RegisterRoutes(globalTestRouter, globalTestContainer.RoutesCfg)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestCreateOrder(t *testing.T) {
	orderPayload := map[string]interface{}{
		"user_id":      "test-user-123",
		"item_ids":     []string{"item1", "item2"},
		"total_amount": 100.50,
	}
	payload, _ := json.Marshal(orderPayload)
	req, _ := http.NewRequest("POST", "/api/v1/orders", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	globalTestRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
	assert.Nil(t, err)
	assert.Contains(t, responseBody, "order_id")
}

func TestGetOrderStatusNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/orders/non-existent-id", nil)
	resp := httptest.NewRecorder()
	globalTestRouter.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestCreateOrderInvalidPayload(t *testing.T) {
	orderPayload := map[string]interface{}{
		"user_id": 123,
	}
	payload, _ := json.Marshal(orderPayload)
	req, _ := http.NewRequest("POST", "/api/v1/orders", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	globalTestRouter.ServeHTTP(resp, req)
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
				"item_ids":     []string{"item1", "item2"},
				"total_amount": 100.50,
			}
			payload, _ := json.Marshal(orderPayload)
			req, _ := http.NewRequest("POST", "/api/v1/orders", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			globalTestRouter.ServeHTTP(resp, req)
			assert.Equal(t, http.StatusOK, resp.Code)
		}(i)
	}
	wg.Wait()
	duration := time.Since(start)
	t.Logf("Completed %v concurrent requests in %s", requestCount, duration)
}

func TestGetMetrics(t *testing.T) {
	time.Sleep(3 * time.Second)
	req, _ := http.NewRequest("GET", "/api/v1/metrics", nil)
	resp := httptest.NewRecorder()
	globalTestRouter.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var metrics map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &metrics)
	assert.Nil(t, err)

	assert.Contains(t, metrics, "total_orders_received")
	assert.Contains(t, metrics, "average_processing_time")
}
