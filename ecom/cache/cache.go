package cache

import (
	"context"
	"log"

	err "ecom.com/errors"
)

var (
	Ctx        = context.Background()
	redisStore map[string]string
)

func InitRedis(addr, password string, dbNum int) {
	redisStore = make(map[string]string)
	log.Println("Connected to Redis")
}

func SetOrderStatus(orderID, status string) error {
	if redisStore == nil {
		return err.ErrUnintializedInstance
	}
	redisStore[orderID] = status
	return nil
}

func GetOrderStatus(orderID string) (string, error) {
	if redisStore == nil {
		return "", err.ErrUnintializedInstance
	}
	val, found := redisStore[orderID]
	if !found {
		return "", err.ErrNotFound
	}
	return val, nil
}
