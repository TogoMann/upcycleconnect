package tests

import (
	"backend/internal/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContainerAPI(t *testing.T) {
	pool := GetPool()
	router := SetupTestRouter(pool)

	token, _ := utils.GenerateJWT(1, "amartin", "client")
	adminToken, _ := utils.GenerateJWT(3, "clefevre", "admin")

	t.Run("CreateLockerAccess", func(t *testing.T) {
		body := map[string]interface{}{
			"item_id": 1,
			"user_id": 1,
		}
		jsonBody, _ := json.Marshal(body)
		
		// Assuming locker ID 1 exists from seeds
		req, _ := http.NewRequest("POST", "/lockers/1/access", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated && rr.Code != http.StatusOK {
			t.Errorf("Expected 201/200, got %v: %s", rr.Code, rr.Body.String())
		}

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)

		if resp["access_code"] == nil || resp["access_code"] == "" {
			t.Errorf("Expected an access code in response, got %v", resp)
		} else {
			fmt.Printf("Generated QR/Access code: %s\n", resp["access_code"])
		}
	})

	t.Run("GetContainers_Admin", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin/conteneurs", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200, got %v: %s", rr.Code, rr.Body.String())
		}
	})
}
