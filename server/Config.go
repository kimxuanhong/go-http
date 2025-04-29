package server

import (
	"os"
)

type Config struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}

func NewConfig() *Config {
	return &Config{
		Host: getEnv("SERVER_HOST", "localhost"),
		Port: getEnv("SERVER_PORT", "8080"),
		Mode: getEnv("GIN_MODE", "debug"),
	}
}

func (c *Config) GetAddr() string {
	return c.Host + ":" + c.Port
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
