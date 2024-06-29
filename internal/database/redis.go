package database

import (
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	// redis_database = os.Getenv("REDIS_DATABASE")
	redis_password = os.Getenv("REDIS_PASSWORD")
	addr           = os.Getenv("ADDR")
	// port       = os.Getenv("DB_PORT")
	// host       = os.Getenv("DB_HOST")
	// schema     = os.Getenv("DB_SCHEMA")
)

func CreateRedisConnection() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redis_password,
		DB:       0,
	})

	return rdb
}
