Features
Order Creation: Create new orders with an initial status of "Pending".
Order Processing: Orders transition through "Pending" → "Processing" → "Completed" states.
Low Latency Reads: Order status is cached in Redis for fast retrieval.
Metrics Reporting: Separate database tracks processing metrics (total processed orders, average processing time) while the orders DB maintains order statuses.
Rotating Logging: Logs are rotated using Lumberjack to prevent unbounded log file growth.
Unit Tests: Comprehensive tests cover API endpoints, database operations, and asynchronous queue processing.

Language: Golang
Web Framework: Gin
Databases: SQLite (separate DBs for orders and metrics)
Cache: Redis (Golang Map)
Queue: In-memory queue with goroutines and channels
Logging: rotating log files
Testing: Go's testing package with Testify

go run main.go

1. Create an Order
Endpoint: POST /orders
Request Body:
json
Copy
Edit
{
  "user_id": "user123",
  "item_ids": "item1,item2",
  "total_amount": 99.99
}
Curl Example:
bash
Copy
Edit
curl -X POST http://localhost:8080/orders \
     -H "Content-Type: application/json" \
     -d '{"user_id": "user123", "item_ids": "item1,item2", "total_amount": 99.99}'
Response:
json
Copy
Edit
{
  "message": "Order created",
  "order_id": "generated-order-id"
}
2. Get Order Status
Endpoint: GET /orders/:order_id
Curl Example:
bash
Copy
Edit
curl http://localhost:8080/orders/<order_id>
Response:
json
Copy
Edit
{
  "order_id": "<order_id>",
  "status": "Pending" // could be "Processing" or "Completed" after processing
}
3. Get Metrics
Endpoint: GET /metrics
Curl Example:
bash
Copy
Edit
curl http://localhost:8080/metrics
Response:
json
Copy
Edit
{
  "total_orders_processed": 10,
  "average_processing_time": 2.1,
  "orders_pending": 0,
  "orders_processing": 0,
  "orders_completed": 10
}
Design Decisions and Trade-offs
Asynchronous Order Processing:
Orders are queued and processed by a worker pool asynchronously, which improves responsiveness. However, it introduces eventual consistency, meaning the order status might not update instantly.

Separate Databases for Orders and Metrics:
Storing orders and metrics in separate databases reduces contention and allows independent scaling of write and read workloads.

Redis for Low Latency:
Order status is cached in Redis to provide quick read access. The cache is updated on order creation and during status transitions. In the event of a cache miss, the system falls back to the orders database.

Rotating Logs:
Using Lumberjack to rotate log files ensures that log files do not grow indefinitely, which is crucial for long-running applications.

Modular Architecture:
The system is divided into distinct modules (handlers, services, models, database, queue, cache, logger) to promote maintainability and ease of testing.

Assumptions
Order Processing Simulation:
Processing is simulated with a fixed delay (using time.Sleep). In a production environment, processing may involve more complex workflows and error handling.

Use of SQLite:
SQLite is used for demonstration and testing purposes. A production system would likely use a more robust database like PostgreSQL or MySQL.

Cache Consistency:
The Redis cache is used for fast reads. It is assumed that eventual consistency is acceptable for order status queries.

Scalability:
The design focuses on modularity and separation of concerns to allow scaling individual components independently as demand increases.

Running Unit Tests
Run the test suite with:

bash
Copy
Edit
go test ./tests -v
The tests cover API endpoints, database operations, and the order processing queue.

