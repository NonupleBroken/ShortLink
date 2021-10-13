package redis

import (
	"ShortLink/logger"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Mutex struct {
	key     string
	value   string
	timeout time.Duration
}

func (m *Mutex) Lock() bool {
	success, err := Client.SetNX(Ctx, m.key, m.value, m.timeout).Result()
	if err != nil {
		logger.S.Error("redis mutex lock error: ", err)
		return false
	}
	return success
}

func (m *Mutex) UnLock() bool {
	deleted, err := Client.Del(Ctx, m.key).Result()
	if err != nil {
		logger.S.Error("redis mutex unlock error: ", err)
		return false
	}
	return deleted > 0
}

func NewMutex(key string, timeout time.Duration) *Mutex {
	redisKey := fmt.Sprintf("_short_link_mutex_%s_", key)
	return &Mutex{
		key:     redisKey,
		value:   uuid.New().String(),
		timeout: timeout,
	}
}
