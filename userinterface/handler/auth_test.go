package handler_test

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"strings"

	"user-service/userinterface/handler"
	"user-service/application/service"
)

func TestAuth_RegisterAndGetUserID(t *testing.T) {
	testToken := "test-token"
	var testUserID uint = 56

	authService := service.NewMockAuthService(testToken, testUserID)

	router := handler.Auth(authService)

	var registerResponse struct {
		Token string `json:"token"`
	}

	/* REGISTER */ {
		registerBody := `{
			"email": "test@example.com",
			"password": "drowssap"
		}`

		req := httptest.NewRequest(
			http.MethodPost,
			"/register/",
			strings.NewReader(registerBody),
		)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("register: expected status %d, got %d", http.StatusOK, rec.Code)
		}

		if err := json.NewDecoder(rec.Body).Decode(&registerResponse); err != nil {
			t.Fatalf("failed to decode register response: %v", err)
		}

		if registerResponse.Token != testToken {
			t.Fatalf(
				"expected token %q, got %q",
				testToken,
				registerResponse.Token,
			)
		}
	}

	/* GET /user/id */ {
		req := httptest.NewRequest(
			http.MethodGet,
			"/user/id",
			nil,
		)
		req.Header.Set(
			"Authorization",
			"Bearer "+registerResponse.Token,
		)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("get user id: expected status %d, got %d. Body: %v",
				http.StatusOK,
				rec.Code,
				rec.Body.String(),
			)
		}

		var userResponse struct {
			ID uint `json:"id"`
		}

		if err := json.NewDecoder(rec.Body).Decode(&userResponse); err != nil {
			t.Fatalf("failed to decode user response: %v", err)
		}

		if userResponse.ID != testUserID {
			t.Fatalf(
				"expected UserID %v, got %v",
				testUserID,
				userResponse.ID,
			)
		}
	}
}
