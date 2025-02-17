# scripts/get_order_status.sh
curl -X GET http://localhost:8081/orders/$1
