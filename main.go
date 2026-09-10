package main

import (
	"user-service/application/service"
	"user-service/infrastructure/cache"
	"user-service/infrastructure/config"
	"user-service/infrastructure/repository"
	"user-service/userinterface"
	"github.com/phuslu/log"
)

func main() {
	config := config.LoadConfig()
	logLevel := log.ParseLevel(config.LogLevel)

	log.DefaultLogger = log.Logger{
		TimeFormat: "15:04:05",
		Level: logLevel,
		Caller: 1,
		Writer: &log.ConsoleWriter{
			ColorOutput: true,
			QuoteString: true,
			EndWithMessage: true,
		},
	}

	userRep := repository.NewUserRepository(config.DatabaseURL)
	userCache := cache.NewUserCache(config.RedisURL)

	authService := service.NewAuthService(userRep, userCache, config.JWTSecret)
	userService := service.NewUserService(userRep, userCache)

	api.Run(config.Port, userService, authService)
}
