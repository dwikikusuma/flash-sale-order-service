package repository

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
)

type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type cacheRepository struct {
	rdb *redis.Client
}

func NewCacheRepository(rdb *redis.Client) CacheRepository {
	return &cacheRepository{
		rdb: rdb,
	}
}

func (r *cacheRepository) Set(ctx context.Context, key string, value interface{}) error {
	err := r.rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *cacheRepository) Get(ctx context.Context, key string) (string, error) {
	value, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", err
	}
	return value, nil
}

func (r *cacheRepository) Delete(ctx context.Context, key string) error {
	err := r.rdb.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
