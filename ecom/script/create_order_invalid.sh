# scripts/create_order_invalid.sh
curl -X POST http://localhost:8081/orders \
-H "Content-Type: application/json" \
-d '{
  "user_id": 123
}'
