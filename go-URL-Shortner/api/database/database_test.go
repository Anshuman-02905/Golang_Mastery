package database

import (
	"os"
	"testing"
	"time"
)

func TestRedisSetAndGet(t *testing.T) {
	os.Setenv("DB_ADDR", "localhost:6379")
	os.Setenv("DB_PASS", "")

	rdb := CreateClient(0)

	key := "test_key"
	value := "hello redis"

	err := rdb.Set(Ctx, key, value, time.Minute).Err()
	if err != nil {
		t.Fatalf("Failed to Set  Key in Redis %v\n", err)
	}
	result, err := rdb.Get(Ctx, key).Result()
	if err != nil {
		t.Fatalf("Failed to Get Key in Redis %v\n", err)
	}

	if result != value {
		t.Fatalf("Expected Value %v  but got %v\n", value, result)

	}
}
