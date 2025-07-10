package ratelimiter

import (
	"context"
	"fmt"
	"time"

	"github.com/joaolima7/rate_limiter-goexpert/internal/domain"
)

type rateLimiter struct {
	storage domain.StorageStrategy
}

func NewRateLimiter(storage domain.StorageStrategy) domain.RateLimiter {
	return &rateLimiter{
		storage: storage,
	}
}

func (rl *rateLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	blocked, err := rl.IsBlocked(ctx, key)
	if err != nil {
		return false, err
	}

	if blocked {
		return false, nil
	}

	count, err := rl.storage.Increment(ctx, fmt.Sprintf("rate:%s", key), window)
	if err != nil {
		return false, err
	}

	return count <= int64(limit), nil
}

func (rl *rateLimiter) Block(ctx context.Context, key string, duration time.Duration) error {
	blockKey := fmt.Sprintf("block:%s", key)
	return rl.storage.Set(ctx, blockKey, "1", duration)
}

func (rl *rateLimiter) IsBlocked(ctx context.Context, key string) (bool, error) {
	blockKey := fmt.Sprintf("block:%s", key)
	return rl.storage.Exists(ctx, blockKey)
}
