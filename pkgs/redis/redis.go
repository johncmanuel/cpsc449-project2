package redis

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

// Initialize Redis client as singleton to prevent multiple
// connections and enable usage via imports
var (
	instance *RedisClient
	once     sync.Once
)

// Default settings for redis
var (
	REDIS_ADDR     = "localhost:6379"
	REDIS_PASSWORD = ""
)

func GetInstance() *RedisClient {
	once.Do(func() {
		rdb := redis.NewClient(&redis.Options{
			Addr:     REDIS_ADDR,
			Password: REDIS_PASSWORD,
			DB:       0,
		})

		// Ping to check connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := rdb.Ping(ctx).Result()
		if err != nil {
			log.Fatalf("Failed to connect to Redis: %v", err)
		}

		instance = &RedisClient{client: rdb}
	})

	return instance
}

func GenerateTupleKey(key1, key2 string) string {
	return fmt.Sprintf("(%s, %s)", key1, key2)
}

// Set a key-value pair with optional expiration
func (r *RedisClient) Set(key string, value interface{}) error {
	defaultExpiration := 2 * time.Minute
	return r.client.Set(context.Background(), key, value, defaultExpiration).Err()
}

// Retrieve a value for a given key
func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(context.Background(), key).Result()
}

// Check if key exists in the cache
func (r *RedisClient) Exists(key string) (bool, error) {
	count, err := r.client.Exists(context.Background(), key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Remove a key
func (r *RedisClient) Delete(key string) error {
	return r.client.Del(context.Background(), key).Err()
}

// Increments the integer value of a key
func (r *RedisClient) Increment(key string) (int64, error) {
	return r.client.Incr(context.Background(), key).Result()
}

// The below hash operations let us store KV pairs (specifically key, string pairs), which can be
// useful for our project.

// Set multiple fields in a hash
func (r *RedisClient) SetHash(key string, fields map[string]interface{}) error {
	return r.client.HMSet(context.Background(), key, fields).Err()
}

// Retrieves all fields of a hash
func (r *RedisClient) GetHash(key string) (map[string]string, error) {
	return r.client.HGetAll(context.Background(), key).Result()
}
