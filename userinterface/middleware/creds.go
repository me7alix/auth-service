package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"user-service/application/dtos"
)

const CredsContextKey string = "credsKey"

func Creds(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var creds dtos.Credentials
		if err := json.Unmarshal(body, &creds); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))
		ctx := context.WithValue(r.Context(), CredsContextKey, creds)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetCredsFromContext(ctx context.Context) dtos.Credentials {
	creds, _ := ctx.Value(CredsContextKey).(dtos.Credentials)
	return creds
}
