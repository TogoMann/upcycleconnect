package payments

import (
	"backend/internal/middlewares"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func getUserId(r *http.Request) (int64, bool) {
	claims, ok := r.Context().Value(middlewares.ClaimsKey).(jwt.MapClaims)
	if !ok {
		return 0, false
	}
	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, false
	}
	return int64(sub), true
}

func (h *Handler) CreateSubscriptionCheckout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId, ok := getUserId(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		PlanId     int64  `json:"plan_id"`
		Siret      string `json:"siret"`
		ReturnPath string `json:"return_path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	url, err := h.service.CreateSubscriptionCheckout(userId, req.PlanId, req.Siret, req.ReturnPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"url": url})
}

func (h *Handler) CreateAdvertisementCheckout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId, ok := getUserId(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := r.PathValue("id")
	adId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid ad id", http.StatusBadRequest)
		return
	}

	url, err := h.service.CreateAdvertisementCheckout(userId, adId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"url": url})
}

func (h *Handler) CreateListingOrderCheckout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId, ok := getUserId(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		ListingId int64 `json:"listing_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	result, err := h.service.CreateListingOrderCheckout(userId, req.ListingId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"free":     result.Free,
		"order_id": result.OrderId,
		"url":      result.URL,
	})
}

func (h *Handler) VerifySession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if _, ok := getUserId(r); !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionId := r.URL.Query().Get("session_id")
	if sessionId == "" {
		http.Error(w, "session_id required", http.StatusBadRequest)
		return
	}

	result, err := h.service.VerifySession(sessionId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"paid": result.Paid,
		"type": result.Type,
	})
}

func (h *Handler) Webhook(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	signatureHeader := r.Header.Get("Stripe-Signature")

	if err := h.service.HandleWebhook(payload, signatureHeader); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
