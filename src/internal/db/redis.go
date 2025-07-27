package db

import (
	"context"
	"course-backend/src/config"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis(addr, password string, db int) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

func InitRedisFromConfig() {
	cfg := config.GetRedisConfig()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.Database,
	})
}

func GetRedisClient() *redis.Client {
	return redisClient
}

func PingRedis(ctx context.Context) error {
	if redisClient == nil {
		return redis.ErrClosed
	}
	return redisClient.Ping(ctx).Err()
}
