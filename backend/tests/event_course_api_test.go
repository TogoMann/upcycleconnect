package tests

import (
	"backend/internal/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestEventCourseAPI(t *testing.T) {
	pool := GetPool()
	router := SetupTestRouter(pool)

	token, _ := utils.GenerateJWT(2, "bdurand", "pro")

	t.Run("CreateEvent_Pro", func(t *testing.T) {
		body := map[string]interface{}{
			"price":      19.99,
			"date":       time.Now().AddDate(0, 0, 7).Format("2006-01-02"),
			"start_time": "14:00:00",
			"end_time":   "16:00:00",
			"location":   "Test Location",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/event/", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated && rr.Code != http.StatusOK {
			t.Errorf("Expected 201/200, got %v: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("CreateCourse_Pro", func(t *testing.T) {
		body := map[string]interface{}{
			"name":         "Test Course",
			"description":  "Course description",
			"max_capacity": 10,
			"price":        49.99,
			"date":         time.Now().AddDate(0, 0, 14).Format("2006-01-02"),
			"start_time":   "10:00:00",
			"end_time":     "12:00:00",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/course", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated && rr.Code != http.StatusOK {
			t.Errorf("Expected 201/200, got %v: %s", rr.Code, rr.Body.String())
		}
	})
}
