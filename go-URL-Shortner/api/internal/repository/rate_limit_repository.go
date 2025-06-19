package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RateLimitRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewRateLimitRepository(client *redis.Client, ctx context.Context) *RateLimitRepository {
	return &RateLimitRepository{
		client: client,
		ctx:    ctx,
	}
}

func (r *RateLimitRepository) GetQouta(ip string) (string, error) {
	return r.client.Get(r.ctx, ip).Result()
}
func (r *RateLimitRepository) SetQuota(ip, quota string, duration time.Duration) error {
	if duration == 0 {
		duration = 30 * time.Minute
	}
	return r.client.Set(r.ctx, ip, quota, duration).Err()
}

func (r *RateLimitRepository) DecrementQuota(ip string) error {
	return r.client.Decr(r.ctx, ip).Err()
}

func (r *RateLimitRepository)Increment(){

}


func (r *RateLimitRepository) GetTTL(ip string) (time.Duration, error) {
	return r.client.TTL(r.ctx, ip).Result()
}
