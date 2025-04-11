package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var (
	ConnectError = errors.New("redis connect error")
	DelErr       = errors.New("value > 1")
	NotFound     = errors.New("value not found")
)

type RedisCache struct {
	*redis.Client
}

type RedisInterface interface {
	Add(key, value string) (string, error)
	Getvalue(key string) (string, error)
	DelKey(key string) error
}

// InitRedis - функция, инициализируем коннект и создаем struct RedisCache.
func InitRedis() (*RedisCache, error) {
	connect := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if err := connect.Ping().Err(); err != nil {
		return nil, fmt.Errorf("%w", ConnectError)
	}

	return &RedisCache{connect}, nil
}

// Add - метод, добавляем данные в redis.
func (r *RedisCache) Add(key, value string) (string, error) {
	resul, err := r.Set(key, value, 100*time.Minute).Result()
	if err != nil {
		return "", err
	}
	return resul, nil
}

// GetValue - метод, ищем значение по ключу
func (r *RedisCache) GetValue(ctx context.Context, key string) error {
	_, err := r.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("%w", NotFound)
		}
		return err
	}
	return nil
}

// DelKey-метод, удаляем ключ из redis.
func (r *RedisCache) DelKey(key string) error {
	ok, err := r.Del(key).Result()
	if err != nil {
		return err
	}
	if ok != 1 {
		return fmt.Errorf("%w", DelErr)
	}
	return nil
}
