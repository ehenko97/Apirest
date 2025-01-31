package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(redisAddr string) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Проверка соединения с Redis
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}

	return &RedisCache{client: rdb}
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	// Redis возвращает строку, но если нужно другое представление, можно сделать преобразование
	*dest.(*string) = val
	return nil
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
