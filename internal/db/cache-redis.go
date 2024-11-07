package db

import (
	"github.com/redis/go-redis/v9"
)

var cacheClient *redis.Client = nil

func InitCacheClient() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	cacheClient = client
}

func GetCacheClient() *redis.Client {
	return cacheClient
}
