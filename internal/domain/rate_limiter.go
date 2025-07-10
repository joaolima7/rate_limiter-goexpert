package domain

import (
	"context"
	"time"
)

type RateLimiter interface {
	Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error)
	Block(ctx context.Context, key string, duration time.Duration) error
	IsBlocked(ctx context.Context, key string) (bool, error)
}

type StorageStrategy interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Increment(ctx context.Context, key string, expiration time.Duration) (int64, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}

type RateLimitRequest struct {
	IP       string
	Token    string
	Endpoint string
}

type RateLimitResponse struct {
	Allowed   bool
	Remaining int64
	ResetTime time.Time
	Reason    string
}
