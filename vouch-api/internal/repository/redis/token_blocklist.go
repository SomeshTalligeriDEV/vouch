package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// TokenBlocklist stores revoked refresh token hashes in Redis.
// Keys expire automatically when the original token would have expired,
// so the set never grows unboundedly.
type TokenBlocklist struct {
	rdb *redis.Client
}

// NewTokenBlocklist constructs a TokenBlocklist backed by the given Redis client.
func NewTokenBlocklist(rdb *redis.Client) *TokenBlocklist {
	return &TokenBlocklist{rdb: rdb}
}

// Block records a token hash as revoked for the given duration.
func (b *TokenBlocklist) Block(ctx context.Context, tokenHash string, exp time.Duration) error {
	return b.rdb.Set(ctx, tokenHash, 1, exp).Err()
}

// IsBlocked reports whether the token hash is in the revocation set.
func (b *TokenBlocklist) IsBlocked(ctx context.Context, tokenHash string) (bool, error) {
	err := b.rdb.Get(ctx, tokenHash).Err()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
