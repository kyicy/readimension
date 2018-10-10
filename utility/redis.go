package utility

import (
	"github.com/go-redis/redis"
)

// RedisClient exports redis client
var RedisClient *redis.Client

// SetUpRedis assign value to RedisClient
func SetUpRedis(client *redis.Client) {
	RedisClient = client
}
