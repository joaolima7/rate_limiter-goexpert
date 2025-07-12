package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	RateLimitIP              int64
	RateLimitToken           int64
	RateLimitDurationSeconds int64
	RedisHost                string
	RedisPort                string
	RedisPassword            string
	RedisDB                  int64
	ServerPort               string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		panic("erro ao carregar variaveis de ambiente: " + err.Error())
	}

	config := &Config{
		RateLimitIP:              getEnvAsInt("RATE_LIMIT_IP"),
		RateLimitToken:           getEnvAsInt("RATE_LIMIT_TOKEN"),
		RateLimitDurationSeconds: getEnvAsInt("RATE_LIMIT_DURATION_SECONDS"),
		RedisHost:                getEnv("REDIS_HOST"),
		RedisPort:                getEnv("REDIS_PORT"),
		RedisPassword:            getEnvOptional("REDIS_PASSWORD"),
		RedisDB:                  getEnvAsInt("REDIS_DB"),
		ServerPort:               getEnv("SERVER_PORT"),
	}

	return config, nil
}

func getEnvAsInt(key string) int64 {
	if value := os.Getenv(key); value != "" {
		valueInt, err := strconv.Atoi(value)
		if err != nil {
			panic("erro ao converter tipo de variavel de ambiente: " + key)
		}

		return int64(valueInt)
	}

	panic("erro ao ler variavel de ambiente: " + key)
}

func getEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	panic("erro ao ler a variavel de ambiente: " + key)
}

func getEnvOptional(key string) string {
	return os.Getenv(key)
}

func (c *Config) GetLimitDuration() time.Duration {
	return time.Duration(c.RateLimitDurationSeconds) * time.Second
}
