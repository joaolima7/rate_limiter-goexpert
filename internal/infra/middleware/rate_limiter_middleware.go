package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/joaolima7/rate_limiter-goexpert/internal/domain"
	"github.com/joaolima7/rate_limiter-goexpert/internal/domain/usecase"
)

type RateLimiterMiddleware struct {
	useCase *usecase.RateLimitUseCase
}

func NewRateLimiterMiddleware(usecase *usecase.RateLimitUseCase) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		useCase: usecase,
	}
}

func (rl *RateLimiterMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := rl.getClientIP(r)
		token := r.Header.Get("API_KEY")

		req := domain.RateLimitRequest{
			IP:       ip,
			Token:    token,
			Endpoint: r.URL.Path,
		}

		response, err := rl.useCase.CheckRateLimit(context.Background(), req)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		if !response.Allowed {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("error: you have reached the maximum number of requests or actions allowed within a certain time frame"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiterMiddleware) getClientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	xri := r.Header.Get("X-Real-Ip")
	if xri != "" {
		return xri
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}
