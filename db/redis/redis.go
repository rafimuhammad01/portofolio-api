package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"os"
)

func Init() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: fmt.Sprintf(os.Getenv("REDIS_PASSWORD")),
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()

	if err != nil {
		logrus.Fatal(fmt.Errorf("error connecting to redis: %w", err))
	}

	logrus.Infoln("Redis Connected Successfully")
	return rdb
}
