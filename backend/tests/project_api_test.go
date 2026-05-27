package tests

import (
	"backend/internal/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProjectAPI(t *testing.T) {
	pool := GetPool()
	router := SetupTestRouter(pool)

	// User 2 (bdurand) is 'pro'
	token, _ := utils.GenerateJWT(2, "bdurand", "pro")

	t.Run("CreateProject", func(t *testing.T) {
		body := map[string]interface{}{
			"listing_id":  1,
			"title":       "Test Project",
			"description": "Project description",
			"final_score": 85,
			"status":      "in progress",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/project", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated && rr.Code != http.StatusOK {
			t.Errorf("Expected 201/200, got %v: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("GetAllProjects", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/project", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200, got %v", rr.Code)
		}
	})
}
