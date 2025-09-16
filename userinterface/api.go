package api

import (
	"fmt"
	"net/http"
	"user-service/application/service"
	"user-service/userinterface/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)


func Run(
	userService service.UserService,
	authService service.AuthService,
) {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Mount("/user", handler.User(userService, authService))
	r.Mount("/auth", handler.Auth(authService))

	fmt.Println("Servier starting on port 8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
