package tests

import (
	"backend/internal/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCartAPI(t *testing.T) {
	pool := GetPool()
	router := SetupTestRouter(pool)

	token, err := utils.GenerateJWT(5, "mdede", "client")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	t.Run("GetCart_Authenticated", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/cart", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", rr.Code)
		}

		var items []interface{}
		if err := json.Unmarshal(rr.Body.Bytes(), &items); err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}
	})

	t.Run("GetCart_Unauthorized", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/cart", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status Unauthorized, got %v", rr.Code)
		}
	})

	t.Run("AddToCart_Listing", func(t *testing.T) {
		body := map[string]interface{}{
			"listing_id": 1,
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/cart", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status Created, got %v", rr.Code)
		}
	})

	t.Run("Checkout", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/cart/checkout", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", rr.Code)
		}
	})
}
