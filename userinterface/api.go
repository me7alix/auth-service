package api

import (
	"fmt"
	"log"
	"net/http"
	"user-service/application/service"
	"user-service/userinterface/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)


func Run(
	port        int,
	userService service.UserService,
	authService service.AuthService,
) {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Mount("/user", handler.User(userService, authService))
	r.Mount("/auth", handler.Auth(authService))

	log.Printf("Server starting on port %v...\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), r); err != nil {
		panic(err)
	}
}
