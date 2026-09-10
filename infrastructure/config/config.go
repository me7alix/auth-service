package config

import (
	"os"
	"fmt"
	"strconv"
)

type Config struct {
	LogLevel    string
	DatabaseURL string
	RedisURL    string
	JWTSecret   string
	Port        int
}

func getEnvStr(key string) string {
	res := os.Getenv(key)
	if res == "" {
		panic(fmt.Errorf("%s is not provided", key))
	}
	return res
}

func getEnvInt(key string) int {
	str := getEnvStr(key)
	res, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(fmt.Errorf("%s parsing error", key))
	}
	return int(res)
}

func LoadConfig() *Config {
	return &Config{
		LogLevel:    getEnvStr("LOG_LEVEL"),
		DatabaseURL: getEnvStr("DATABASE_URL"),
		RedisURL:    getEnvStr("REDIS_URL"),
		JWTSecret:   getEnvStr("JWT_SECRET"),
		Port:        getEnvInt("PORT"),
	}
}
