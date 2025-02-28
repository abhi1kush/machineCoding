package cache

import (
	"context"
	"log"
	"sync"

	err "ecom.com/errors"
)

var (
	Ctx               = context.Background()
	redisStore        map[string]string
	redisStoreRWMutex *sync.RWMutex
)

func InitRedis(addr, password string, dbNum int) {
	redisStore = make(map[string]string)
	redisStoreRWMutex = &sync.RWMutex{}
	log.Println("Connected to Redis")
}

func SetOrderStatus(orderID, status string) error {
	if redisStore == nil {
		return err.ErrUnintializedInstance
	}
	redisStoreRWMutex.Lock()
	redisStore[orderID] = status
	redisStoreRWMutex.Unlock()
	return nil
}

func GetOrderStatus(orderID string) (string, error) {
	if redisStore == nil {
		return "", err.ErrUnintializedInstance
	}
	redisStoreRWMutex.RLock()
	val, found := redisStore[orderID]
	redisStoreRWMutex.RUnlock()
	if !found {
		return "", err.ErrNotFound
	}
	return val, nil
}
