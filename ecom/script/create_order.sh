# scripts/create_order.sh
curl -X POST http://localhost:8081/orders \
-H "Content-Type: application/json" \
-d '{
  "user_id": "test-user-123",
  "item_ids": "item1,item2",
  "total_amount": 100.50
}'
