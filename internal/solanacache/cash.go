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
	}

	// Option is a function that configures the SolanaClientCacheWrapper
	Option func(*SolanaClientCacheWrapper)
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
		if err := c.cache.Get(ctx, base58MintAddr, &result); err == nil {
			return &result, nil
		}
	}

	metadata, err := c.Client.GetTokenMetadata(ctx, base58MintAddr)
	if err != nil {
		return nil, err
	}

	if metadata == nil {
		return nil, fmt.Errorf("token metadata is nil")
	}

	// cache the result, ignore error, bc it is not critical
	c.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   base58MintAddr,
		Value: metadata,
		TTL:   c.ttl,
	})

	return metadata, nil
}
