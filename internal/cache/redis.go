package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aeilang/urlshortener/config"
	"github.com/aeilang/urlshortener/internal/repo"
	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	cleint *redis.Client
}

func NewReisClient(cfg config.RedisConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{
		cleint: client,
	}, nil
}

func (r *RedisCache) SetURL(ctx context.Context, url repo.Url) error {
	data, err := json.Marshal(url)
	if err != nil {
		return err
	}
	if url.ExpiredAt.Before(time.Now()) {
		return nil
	}

	return r.cleint.Set(ctx, url.ShortCode, data, time.Until(url.ExpiredAt)).Err()
}

func (r *RedisCache) GetURLByShortCode(ctx context.Context, shortCode string) (*repo.Url, error) {
	data, err := r.cleint.Get(ctx, shortCode).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var url repo.Url
	if err := json.Unmarshal(data, &url); err != nil {
		return nil, err
	}

	return &url, nil
}

func (r *RedisCache) Close() error {
	return r.cleint.Close()
}
