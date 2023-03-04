package cache

import (
	"log"

	"github.com/redis/go-redis/v9"
)

type Connection struct {
	redisClient *redis.Client
}

type Provider interface {
	GetRedisCache() *redis.Client
}

func ConnectRedis(redisUrl string) *Connection {
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		log.Fatalf("Redis URL is invalid: %v\n", err)
	}

	return &Connection{
		redisClient: redis.NewClient(opt),
	}
}

func (c *Connection) CloseRedisClient() {
	if err := c.redisClient.Close(); err != nil {
		log.Fatalf("Failed to close Redis! => %v\n", err)
	}
}

func (c *Connection) GetRedisCache() *redis.Client {
	return c.redisClient
}
