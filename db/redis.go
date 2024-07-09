package db

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis() {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	redisClient = client
}

func GetRedisClient() *redis.Client {
	return redisClient
}

func ClearRedisClient() {
	status := GetRedisClient().FlushDB(context.Background())
	fmt.Println(status.Args()...)
}
