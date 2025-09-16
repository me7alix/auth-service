package main

import (
	"user-service/application/service"
	"user-service/infrastructure/cache"
	"user-service/infrastructure/config"
	"user-service/infrastructure/repository"
	"user-service/userinterface"
)

func main() {
	config := config.LoadConfig()

	userRep := repository.NewUserRepository(config.DatabaseURL)
	userCache := cache.NewUserCache(config.RedisURL)

	authService := service.NewAuthService(userRep, userCache, config.JWTSecret)
	userService := service.NewUserService(userRep)

	api.Run(userService, authService)
}
