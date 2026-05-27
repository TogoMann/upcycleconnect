package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthAPI(t *testing.T) {
	pool := GetPool()
	router := SetupTestRouter(pool)

	uniqueID := time.Now().UnixNano()
	email := fmt.Sprintf("test_%d@example.com", uniqueID)
	password := "password123"
	username := fmt.Sprintf("user_%d", uniqueID)

	t.Run("Register_NewUser", func(t *testing.T) {
		body := map[string]string{
			"username":   username,
			"email":      email,
			"password":   password,
			"first_name": "Test",
			"last_name":  "User",
			"role":       "client",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/register/", bytes.NewBuffer(jsonBody))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated && rr.Code != http.StatusOK {
			t.Errorf("Expected status Created/OK, got %v: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("Login_ValidCredentials", func(t *testing.T) {
		body := map[string]string{
			"username": username,
			"password": password,
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/login/", bytes.NewBuffer(jsonBody))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v: %s", rr.Code, rr.Body.String())
		}

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp["token"] == "" {
			t.Error("Expected token in response")
		}
	})

	t.Run("Login_InvalidCredentials", func(t *testing.T) {
		body := map[string]string{
			"email":    email,
			"password": "wrongpassword",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/login/", bytes.NewBuffer(jsonBody))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code == http.StatusOK {
			t.Error("Expected failure for invalid credentials")
		}
	})
}
