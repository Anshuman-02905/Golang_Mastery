package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type URLrepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewUrlRepository(client *redis.Client, ctx context.Context) *URLrepository {
	return &URLrepository{client: client, ctx: ctx}
}
func (r URLrepository) Exists(id string) (bool, string, error) {
	val, err := r.client.Get(r.ctx, id).Result()
	if err == redis.Nil {
		return false, "", nil
	}
	return true, val, err
}

func (r *URLrepository)Save(id ,url string,exipry time.Duration)error{
	return r.client.Set(r.ctx,id,url,exipry).Err()
}
