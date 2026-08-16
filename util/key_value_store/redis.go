package key_value_store

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisWrapper struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedisConnection(redisHost, redisUser, redisPassword string) (*RedisWrapper, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Username: redisUser,
		Password: redisPassword,
		DB:       0,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis: %w", err)
		return nil, fmt.Errorf("Failed to connect to Redis: %w", err)
	}

	log.Println("Connected to Redis")

	return &RedisWrapper{Client: client, Ctx: ctx}, nil
}

func (r *RedisWrapper) Set(key string, value interface{}, ttl time.Duration) error {
	jsonVal, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return r.Client.Set(r.Ctx, key, jsonVal, ttl).Err()
}

func (r *RedisWrapper) Get(key string, dest interface{}) error {
	val, err := r.Client.Get(r.Ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisWrapper) Delete(key string) error {
	return r.Client.Del(r.Ctx, key).Err()
}

func (r *RedisWrapper) SetWithTTL(key string, value string, ttl time.Duration) error {
	return r.Client.Set(r.Ctx, key, value, ttl).Err()
}

func (r *RedisWrapper) Expire(key string, ttl time.Duration) error {
	return r.Client.Expire(r.Ctx, key, ttl).Err()
}

func (r *RedisWrapper) AddToSet(key string, value string) error {
	return r.Client.SAdd(r.Ctx, key, value).Err()
}

func (r *RedisWrapper) RemoveFromSet(key string, value string) error {
	return r.Client.SRem(r.Ctx, key, value).Err()
}
