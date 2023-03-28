package solanacache

import (
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

// WithCacheTTL sets the cache ttl
func WithCacheTTL(ttl time.Duration) Option {
	return func(c *SolanaClientCacheWrapper) {
		c.ttl = ttl
	}
}

// WithCacheClient sets the cache client
func WithCacheClient(cache *cache.Cache) Option {
	return func(c *SolanaClientCacheWrapper) {
		c.cache = cache
	}
}

// WithRedisClient inits cache client with given redis connection.
func WithRedisClient(redisClient *redis.Client) Option {
	cacheClient := cache.New(&cache.Options{
		Redis:      redisClient,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
	return func(c *SolanaClientCacheWrapper) {
		c.cache = cacheClient
	}
}

// WithRedisConnectionURL init redis connection and cache client with it.
func WithRedisConnectionURL(redisConnURL string) Option {
	redisOpt, err := redis.ParseURL(redisConnURL)
	if err != nil {
		panic("Failed to parse redis connection url")
	}
	cacheClient := cache.New(&cache.Options{
		Redis:      redis.NewClient(redisOpt),
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
	return func(c *SolanaClientCacheWrapper) {
		c.cache = cacheClient
	}
}
