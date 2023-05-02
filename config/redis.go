package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisCli *redis.Client

func InitRedis() {
	if redisCli == nil {
		redisHost := os.Getenv("REDIS-HOST")
		if redisHost == "" {
			redisHost = "127.0.0.1"
		}
		redisCli = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:6379", redisHost),
			Password: "", // 没有密码，默认值
			DB:       0,  // 默认DB 0
		})
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := redisCli.Ping(ctx).Err(); err != nil {
			log.Fatal(err)
		}
	}
}

func GetRedisClient() *redis.Client {
	return redisCli
}
