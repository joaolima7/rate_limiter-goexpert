package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaolima7/rate_limiter-goexpert/internal/config"
	"github.com/joaolima7/rate_limiter-goexpert/internal/domain/usecase"
	"github.com/joaolima7/rate_limiter-goexpert/internal/infra/middleware"
	"github.com/joaolima7/rate_limiter-goexpert/internal/infra/storage"
	"github.com/joaolima7/rate_limiter-goexpert/pkg/ratelimiter"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	storage, err := storage.NewRedisStorage(cfg)
	if err != nil {
		log.Fatal("failed to initialize Redis storage: ", err)
	}

	rateLimiter := ratelimiter.NewRateLimiter(storage)

	rateLimitUseCase := usecase.NewRateLimiterUseCase(rateLimiter, cfg)

	rateLimitMiddleware := middleware.NewRateLimiterMiddleware(rateLimitUseCase)

	router := mux.NewRouter()

	router.Use(rateLimitMiddleware.Handler)

	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/test", testHandler).Methods("GET")

	log.Printf("server starting on port: %s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, router))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "message: Rate Limiter Service is running")
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "message: Test endpoint, status: success")
}
