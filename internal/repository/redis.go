package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
)

func NewRedis(addr, password, dbStr string) *redis.Client {
	db, err := strconv.Atoi(dbStr)
	if err != nil {
		db = 0
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Redis ping failed: %v", err)
	}
	return rdb
}
