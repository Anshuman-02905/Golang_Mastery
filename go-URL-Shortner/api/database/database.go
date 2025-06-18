package database

import (
	"context"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

func CreateClient(dbNo int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDR"),
		Password: os.Getenv("DB_PASS"),
		DB:       dbNo,
	})

	return rdb

}

var (
	onceDB0 sync.Once
	onceDB1 sync.Once
	rdb0    *redis.Client
	rdb1    *redis.Client
)

func GetClient(dbno int) *redis.Client {

	switch dbno {
	case 0:
		onceDB0.Do(func() {
			rdb0 = redis.NewClient(&redis.Options{
				Addr:     os.Getenv("DB_ADDR"),
				Password: os.Getenv("DB_PASS"),
				DB:       0,
			})
		})
		return rdb0
	case 1:
		onceDB1.Do(func() {
			rdb1 = redis.NewClient(&redis.Options{
				Addr:     os.Getenv("DB_ADDR"),
				Password: os.Getenv("DB_PASS"),
				DB:       1,
			})
		})
		return rdb1
	default:
		panic("Unsupported DB Number")
	}

}
