package solanacache

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitrymomot/solana/client"
	"github.com/dmitrymomot/solana/token_metadata"
	"github.com/go-redis/cache/v8"
)

type (
	// SolanaClientCacheWrapper is a wrapper for solana client
	SolanaClientCacheWrapper struct {
		*client.Client
		cache *cache.Cache
		ttl   time.Duration
		log   logger
	}

	// Option is a function that configures the SolanaClientCacheWrapper
	Option func(*SolanaClientCacheWrapper)

	// logger is a logger interface
	logger interface {
		Printf(format string, v ...interface{})
		Debugf(format string, v ...interface{})
		Errorf(format string, v ...interface{})
	}
)

// NewSolanaClientCacheWrapper creates a new instance of SolanaClientCacheWrapper
func NewSolanaClientCacheWrapper(c *client.Client, opts ...Option) *SolanaClientCacheWrapper {
	wrapper := &SolanaClientCacheWrapper{
		Client: c,
		ttl:    time.Hour,
	}

	for _, opt := range opts {
		opt(wrapper)
	}

	return wrapper
}

// GetTokenMetadata returns token metadata.
// It is used to cache the result of the solana client method.
func (c *SolanaClientCacheWrapper) GetTokenMetadata(ctx context.Context, base58MintAddr string) (*token_metadata.Metadata, error) {
	var result token_metadata.Metadata
	if c.cache.Exists(ctx, base58MintAddr) {
		c.log.Debugf("cache exists for %s", base58MintAddr)
		if err := c.cache.Get(ctx, base58MintAddr, &result); err == nil && result.Mint != "" {
			c.log.Debugf("cache hit for %s: metadta: %+v", base58MintAddr, result)
			return &result, nil
		}
	}

	metadata, err := c.Client.GetTokenMetadata(ctx, base58MintAddr)
	if err != nil {
		return nil, err
	}

	if metadata == nil {
		c.log.Errorf("token metadata is nil for %s", base58MintAddr)
		return nil, fmt.Errorf("token metadata is nil")
	}

	// cache the result, ignore error, bc it is not critical
	if err := c.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   base58MintAddr,
		Value: metadata,
		TTL:   c.ttl,
	}); err != nil {
		c.log.Errorf("failed to cache token metadata for %s: %s", base58MintAddr, err.Error())
	}

	return metadata, nil
}
