package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/joaolima7/rate_limiter-goexpert/internal/config"
	"github.com/joaolima7/rate_limiter-goexpert/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRateLimiter struct {
	mock.Mock
}

func (m *MockRateLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	args := m.Called(ctx, key, limit, window)
	return args.Bool(0), args.Error(1)
}

func (m *MockRateLimiter) Block(ctx context.Context, key string, duration time.Duration) error {
	args := m.Called(ctx, key, duration)
	return args.Error(0)
}

func (m *MockRateLimiter) IsBlocked(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func TestRateLimitUseCase_CheckRateLimit_Token(t *testing.T) {
	mockLimiter := new(MockRateLimiter)
	cfg := &config.Config{
		RateLimitToken:           100,
		RateLimitIP:              10,
		RateLimitDurationSeconds: 300,
	}

	useCase := NewRateLimiterUseCase(mockLimiter, cfg)

	mockLimiter.On("Allow", mock.Anything, "token:test-token", 100, time.Second).Return(true, nil)

	req := domain.RateLimitRequest{
		IP:    "192.168.1.1",
		Token: "test-token",
	}

	response, err := useCase.CheckRateLimit(context.Background(), req)

	assert.NoError(t, err)
	assert.True(t, response.Allowed)
	mockLimiter.AssertExpectations(t)
}
