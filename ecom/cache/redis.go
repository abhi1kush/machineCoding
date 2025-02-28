package cache

import (
	"context"
	"sync"

	err "ecom.com/errors"
)

type Redis struct {
	Ctx               context.Context
	redisStore        map[string]string
	redisStoreRWMutex *sync.RWMutex
}

func NewRedis(addr, password string, dbNum int) CacheI {
	return &Redis{
		Ctx:               context.Background(),
		redisStore:        map[string]string{},
		redisStoreRWMutex: &sync.RWMutex{},
	}
}

func (r *Redis) SetOrderStatus(orderID, status string) error {
	if r.redisStore == nil {
		return err.ErrUnintializedInstance
	}
	r.redisStoreRWMutex.Lock()
	r.redisStore[orderID] = status
	r.redisStoreRWMutex.Unlock()
	return nil
}

func (r *Redis) GetOrderStatus(orderID string) (string, error) {
	if r.redisStore == nil {
		return "", err.ErrUnintializedInstance
	}
	r.redisStoreRWMutex.RLock()
	val, found := r.redisStore[orderID]
	r.redisStoreRWMutex.RUnlock()
	if !found {
		return "", err.ErrNotFound
	}
	return val, nil
}
