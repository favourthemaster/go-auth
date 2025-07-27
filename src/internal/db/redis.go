package db

import (
	"authentication/src/config"
	"context"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// InitRedis initializes the Redis client with the given address, password, and database number.
func InitRedis(addr, password string, db int) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

// InitRedisFromConfig initializes the Redis client using configuration values.
func InitRedisFromConfig() {
	cfg := config.GetRedisConfig()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.Database,
	})
}

// GetRedisClient returns the Redis client instance.
func GetRedisClient() *redis.Client {
	return redisClient
}

// PingRedis checks the connection to the Redis server.
func PingRedis(ctx context.Context) error {
	if redisClient == nil {
		return redis.ErrClosed
	}
	return redisClient.Ping(ctx).Err()
}
