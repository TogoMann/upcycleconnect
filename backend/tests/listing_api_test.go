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

func TestListingAPI(t *testing.T) {
	pool := GetPool()
	router := SetupTestRouter(pool)

	token, _ := utils.GenerateJWT(5, "mdede", "client")
	adminToken, _ := utils.GenerateJWT(3, "clefevre", "admin")

	var createdListingId int64

	t.Run("CreateListing", func(t *testing.T) {
		body := map[string]interface{}{
			"name":        "Test Listing",
			"description": "Description",
			"price":       42.50,
			"category":    "Mobilier",
			"city_id":      1,
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/listing/", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected 201, got %v: %s", rr.Code, rr.Body.String())
		}

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		
		rawId := resp["id"]
		if idMap, ok := rawId.(map[string]interface{}); ok {
			createdListingId = int64(idMap["Int64"].(float64))
		} else {
			createdListingId = int64(rawId.(float64))
		}
	})

	t.Run("GetAllApproved", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/listing", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200, got %v", rr.Code)
		}
	})

	t.Run("ApproveListing_Admin", func(t *testing.T) {
		req, _ := http.NewRequest("PATCH", fmt.Sprintf("/listing/%d/approve", createdListingId), nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200, got %v: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("DeleteListing", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/listing/%d", createdListingId), nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200, got %v: %s", rr.Code, rr.Body.String())
		}
	})
}
