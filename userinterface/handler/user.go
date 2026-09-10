package handler

import (
	"encoding/json"
	"net/http"
	"user-service/application/service"
	"user-service/userinterface/middleware"

	"github.com/go-chi/chi/v5"
)

func User(userService service.UserService, authService service.AuthService) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Auth(authService))

	r.Patch("/nickname", func (w http.ResponseWriter, r *http.Request) {
		user := middleware.GetUserFromContext(r.Context())

		var data map[string]any
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		err = userService.UpdateNickname(user.ID, data["nickname"].(string))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	r.Patch("/password", func (w http.ResponseWriter, r *http.Request) {
		user := middleware.GetUserFromContext(r.Context())

		var data map[string]any
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		err = userService.UpdatePassword(user.ID, data["password"].(string))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	r.Patch("/pfp", func (w http.ResponseWriter, r *http.Request) {
		user := middleware.GetUserFromContext(r.Context())

		var data map[string]any
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		err = userService.UpdateProfilePicture(user.ID, data["pfp"].(string))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	return r
}
