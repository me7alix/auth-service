package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"user-service/application/service"
	"user-service/userinterface/middleware"

	"github.com/go-chi/chi/v5"
)

func Auth(authService service.AuthService) http.Handler {
	r := chi.NewRouter()

	r.Route("/user", func (r chi.Router) {
		r.Use(middleware.Auth(authService))

		r.Get("/", func (w http.ResponseWriter, r *http.Request) {
			user, _ := middleware.GetUserFromContext(r.Context())

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
		})

		r.Get("/id", func (w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}

			const prefix = "Bearer "
			if !strings.HasPrefix(authHeader, prefix) {
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, prefix)
			if token == "" {
				http.Error(w, "Token is empty", http.StatusUnauthorized)
				return
			}

			userID, err := authService.GetID(token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"id": %d}`, userID)
		})
	})

	r.Route("/login", func (r chi.Router) {
		r.Use(middleware.Creds)

		r.Post("/", func (w http.ResponseWriter, r *http.Request) {
			creds, _ := middleware.GetCredsFromContext(r.Context())

			token, err := authService.Login(creds)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			fmt.Fprintf(w, `{"token": "%s"}`, token)
		})
	})

	r.Route("/register", func (r chi.Router) {
		r.Use(middleware.Creds)

		r.Post("/", func (w http.ResponseWriter, r *http.Request) {
			creds, _ := middleware.GetCredsFromContext(r.Context())

			token, err := authService.Register(creds)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			fmt.Fprintf(w, `{"token": "%s"}`, token)
		})
	})

	return r
}
