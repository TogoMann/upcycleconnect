package tests

import (
	"backend/internal/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChatAPI(t *testing.T) {
	pool := GetPool()
	router := SetupTestRouter(pool)

	token1, _ := utils.GenerateJWT(1, "amartin", "client")

	t.Run("SendMessage_CreateConversation", func(t *testing.T) {
		body := map[string]interface{}{
			"listing_id":   1,
			"content":      "Hello, I am interested!",
			"message_type": "text",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/chat/messages", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer "+token1)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated && rr.Code != http.StatusOK {
			t.Errorf("Expected 201/200, got %v: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("GetConversations", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/chat/conversations", nil)
		req.Header.Set("Authorization", "Bearer "+token1)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200, got %v: %s", rr.Code, rr.Body.String())
		}

		var convs []interface{}
		json.Unmarshal(rr.Body.Bytes(), &convs)
		if len(convs) == 0 {
			t.Log("Warning: No conversations found for user 1")
		}
	})

	t.Run("SendPriceProposal", func(t *testing.T) {
		price := 120.0
		body := map[string]interface{}{
			"listing_id":     1,
			"content":        "I propose 120",
			"message_type":   "price_proposal",
			"proposed_price": price,
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/chat/messages", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer "+token1)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated && rr.Code != http.StatusOK {
			t.Errorf("Expected 201/200, got %v: %s", rr.Code, rr.Body.String())
		}
	})
}
