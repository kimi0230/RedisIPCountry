package connect

import (
	"RedisIPCountry/config"
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func ConnectRedis() *redis.Client {
	fmt.Println("Redis Server : ", config.Addr)
	conn := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	if _, err := conn.Ping(ctx).Result(); err != nil {
		log.Fatalf("Connect to redis client failed, err: %v\n", err)
	}
	return conn
}
