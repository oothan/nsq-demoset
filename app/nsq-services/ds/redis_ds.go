package ds

import (
	"fmt"
	"github.com/go-redis/redis/v9"
	logger "nsq-demoset/app/_applib"
	"os"
)

var RDB *redis.Client

func LoadRDB() (*redis.Client, error) {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	user := os.Getenv("REDIS_USER")
	pass := os.Getenv("REDIS_PASS")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Username: user,
		Password: pass,
		DB:       0,
	})

	RDB = rdb

	logger.Sugar.Info("Successfully connected to redis")

	return rdb, nil
}
