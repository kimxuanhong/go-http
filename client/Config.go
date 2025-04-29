package client

import (
	"os"
	"strconv" // <- Thêm dòng này
	"time"
)

type Config struct {
	BaseURL    string            `yaml:"base_url"`
	Timeout    time.Duration     `yaml:"timeout"`
	RetryCount int               `yaml:"retry_count"`
	RetryWait  time.Duration     `yaml:"retry_wait"`
	Headers    map[string]string `yaml:"headers"`
}

func NewConfig() *Config {
	timeout, _ := time.ParseDuration(getEnv("CLIENT_TIMEOUT", "30s"))
	retryWait, _ := time.ParseDuration(getEnv("CLIENT_RETRY_WAIT", "1s"))
	retryCount, _ := strconv.Atoi(getEnv("CLIENT_RETRY_COUNT", "3"))

	return &Config{
		BaseURL:    getEnv("CLIENT_BASE_URL", ""),
		Timeout:    timeout,
		RetryCount: retryCount,
		RetryWait:  retryWait,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
