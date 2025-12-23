package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/masudur-rahman/expense-tracker-bot/configs"

	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Cacher struct {
	rc *redis.Client
	mc *cache.Cache
}

var cacher = Cacher{}

func Init(config configs.CacheConfig) {
	switch config.Type {
	case configs.CacheMap:
		cacher.mc = cache.New(-1, -1)
	case configs.CacheRedis:
		cacher.rc = redis.NewClient(&redis.Options{
			Addr:       fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port),
			ClientName: "expense-tracker",
			Username:   config.Redis.User,
			Password:   config.Redis.Password,
			DB:         0,
		})
		_, err := cacher.rc.Ping(ctx).Result()
		if err != nil {
			log.Fatalf("Could not connect to Redis: %v", err)
		}
		fmt.Println("Connected to Redis!")
	}

}

// SetCache stores a key-value pair in Redis with a specific expiration time
func SetCache(key string, value string, expiration time.Duration) (err error) {
	if cacher.mc != nil {
		cacher.mc.Set(key, value, expiration)
	} else if cacher.rc != nil {
		err = cacher.rc.Set(ctx, key, value, expiration).Err()
	}
	if err != nil {
		return fmt.Errorf("failed to set cache key %s: %w", key, err)
	}
	return nil
}

// GetCache retrieves a value from Redis
func GetCache(key string) (string, bool) {
	if cacher.mc != nil {
		if val, has := cacher.mc.Get(key); has {
			return val.(string), has
		}
		return "", false
	} else if cacher.rc != nil {
		val, err := cacher.rc.Get(ctx, key).Result()
		if err != nil {
			return "", false
		}
		return val, true
	}
	return "", false
}
