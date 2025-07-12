package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/joaolima7/rate_limiter-goexpert/internal/config"
	"github.com/joaolima7/rate_limiter-goexpert/internal/domain"
)

type RateLimitUseCase struct {
	rateLimiter domain.RateLimiter
	config      *config.Config
}

func NewRateLimiterUseCase(rateLimiter domain.RateLimiter, config *config.Config) *RateLimitUseCase {
	return &RateLimitUseCase{
		rateLimiter: rateLimiter,
		config:      config,
	}
}

func (uc *RateLimitUseCase) CheckRateLimit(ctx context.Context, req domain.RateLimitRequest) (*domain.RateLimitResponse, error) {
	window := time.Second

	if req.Token != "" {
		key := fmt.Sprintf("token:%s", req.Token)
		allowed, err := uc.rateLimiter.Allow(ctx, key, int(uc.config.RateLimitToken), window)
		if err != nil {
			return nil, err
		}

		if !allowed {

			if err := uc.rateLimiter.Block(ctx, key, uc.config.GetLimitDuration()); err != nil {
				return nil, err
			}

			return &domain.RateLimitResponse{
				Allowed: false,
				Reason:  "Token rate limit exceeded",
			}, nil
		}

		return &domain.RateLimitResponse{
			Allowed: true,
		}, nil
	}

	key := fmt.Sprintf("ip:%s", req.IP)
	allowed, err := uc.rateLimiter.Allow(ctx, key, int(uc.config.RateLimitIP), window)
	if err != nil {
		return nil, err
	}

	if !allowed {

		if err := uc.rateLimiter.Block(ctx, key, uc.config.GetLimitDuration()); err != nil {
			return nil, err
		}

		return &domain.RateLimitResponse{
			Allowed: false,
			Reason:  "IP rate limit exceeded",
		}, nil
	}

	return &domain.RateLimitResponse{
		Allowed: true,
	}, nil

}
