package solanacache

import (
	"context"
	"encoding/json"
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
	// try to get metadata from cache
	metadata, err := c.getMetadataFromCache(ctx, base58MintAddr)
	if err != nil {
		c.log.Errorf("failed to get token metadata from cache for %s: %s", base58MintAddr, err.Error())
	}
	if metadata != nil && metadata.Data != nil && metadata.Mint == base58MintAddr {
		return metadata, nil
	}

	metadata, err = c.Client.GetTokenMetadata(ctx, base58MintAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get token metadata from solana: %w", err)
	}

	if metadata == nil {
		c.log.Errorf("token metadata is nil for %s", base58MintAddr)
		return nil, fmt.Errorf("token metadata is empty")
	}

	if err := c.cacheMetadata(ctx, base58MintAddr, metadata); err != nil {
		c.log.Errorf("failed to cache token metadata for %s: %s", base58MintAddr, err.Error())
	}

	return metadata, nil
}

// marshal metadata to json string and cache it
func (c *SolanaClientCacheWrapper) cacheMetadata(ctx context.Context, base58MintAddr string, metadata *token_metadata.Metadata) error {
	if metadata == nil {
		c.log.Errorf("token metadata is nil for %s", base58MintAddr)
		return fmt.Errorf("token metadata is nil")
	}

	val, err := json.Marshal(metadata)
	if err != nil {
		c.log.Errorf("failed to marshal token metadata for %s: %s", base58MintAddr, err.Error())
		return err
	}

	// cache the result, ignore error, bc it is not critical
	if err := c.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   base58MintAddr,
		Value: string(val),
		TTL:   c.ttl,
	}); err != nil {
		c.log.Errorf("failed to cache token metadata for %s: %s", base58MintAddr, err.Error())
		return err
	}

	return nil
}

// get metadata from cache
func (c *SolanaClientCacheWrapper) getMetadataFromCache(ctx context.Context, base58MintAddr string) (*token_metadata.Metadata, error) {
	if !c.cache.Exists(ctx, base58MintAddr) {
		c.log.Debugf("cache does not exist for %s", base58MintAddr)
		return nil, nil
	} else {
		c.log.Debugf("cache exists for %s", base58MintAddr)
	}

	var data string
	if err := c.cache.Get(ctx, base58MintAddr, &data); err != nil {
		return nil, fmt.Errorf("failed to get token metadata from cache for %s: %w", base58MintAddr, err)
	}

	var result token_metadata.Metadata
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token metadata for %s: %w", base58MintAddr, err)
	}

	return &result, nil
}
