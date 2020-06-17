package cache

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisConnection struct {
	host     string
	password string
	db       int
	expires  time.Duration
}

// RedisClient ...
type RedisClient struct {
	*redis.Client
}

// NewRedisConnection ...
func NewRedisConnection(host string, password string, db int, expires time.Duration) Cacher {
	return &redisConnection{
		host:     host,
		password: password,
		db:       db,
		expires:  expires,
	}
}

var once sync.Once
var redisClient *RedisClient

func (cache *redisConnection) getClient() *RedisClient {
	once.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     cache.host,
			Password: cache.password,
			DB:       cache.db,
			PoolSize: 10,
		})
		redisClient = &RedisClient{client}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Unable to connect to redis %v", err)
	}

	return redisClient
}

// Get ...
func (cache *redisConnection) Get(ctx context.Context, key string, src *Cachable) error {
	client := cache.getClient()

	value, err := client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}

	err = json.Unmarshal([]byte(value), &src)
	if err != nil {
		return err
	}

	return nil
}

// Set ...
func (cache *redisConnection) Set(ctx context.Context, key string, value *Cachable) error {
	client := cache.getClient()

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = client.Set(ctx, key, data, cache.expires*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}
