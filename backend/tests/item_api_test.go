package tests

import (
	"backend/internal/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestItemAPI(t *testing.T) {
	pool := GetPool()
	router := SetupTestRouter(pool)

	token, _ := utils.GenerateJWT(1, "amartin", "client")

	t.Run("CreateItem", func(t *testing.T) {
		body := map[string]interface{}{
			"material_type": "Bois",
			"physical_state": "bon etat",
			"status":        "deposited",
			"site_id":       1,
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/items/", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated && rr.Code != http.StatusOK {
			t.Errorf("Expected 201/200, got %v: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("GetMyItems", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/items/me", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200, got %v: %s", rr.Code, rr.Body.String())
		}
	})
}
