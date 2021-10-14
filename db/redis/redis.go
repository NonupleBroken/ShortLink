package redis

import (
	"ShortLink/config"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var Client *redis.Client
var Ctx = context.Background()
var Nil = redis.Nil

func InitRedis(cfg *config.RedisConfig) error {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	_, err := Client.Ping(Ctx).Result()
	return err
}

