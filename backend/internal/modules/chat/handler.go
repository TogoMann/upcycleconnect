package chat

import (
	"backend/internal/middlewares"
	"encoding/json"
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

func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	msg, err := h.service.SendMessage(userId, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

func (h *Handler) GetConversationMessages(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	convIdStr := r.PathValue("id")
	convId, err := strconv.ParseInt(convIdStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid conversation ID", http.StatusBadRequest)
		return
	}

	messages, err := h.service.GetConversationMessages(userId, convId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (h *Handler) EditMessage(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	msgIdStr := r.PathValue("id")
	msgId, err := strconv.ParseInt(msgIdStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid message ID", http.StatusBadRequest)
		return
	}

	var req EditMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.service.EditMessage(userId, msgId, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) HandleProposal(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	msgIdStr := r.PathValue("id")
	msgId, err := strconv.ParseInt(msgIdStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid message ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Accept bool `json:"accept"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.service.HandleProposal(userId, msgId, req.Accept); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetUserConversations(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	convs, err := h.service.GetUserConversations(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(convs)
}

func (h *Handler) AdminGetConversations(w http.ResponseWriter, r *http.Request) {
	convs, err := h.service.AdminGetConversations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(convs)
}

func (h *Handler) AdminGetConversationDetails(w http.ResponseWriter, r *http.Request) {
	convIdStr := r.PathValue("id")
	convId, err := strconv.ParseInt(convIdStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid conversation ID", http.StatusBadRequest)
		return
	}

	messages, history, err := h.service.AdminGetConversationDetails(convId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"messages": messages,
		"history":  history,
	})
}

func getUserIdFromContext(r *http.Request) (int64, error) {
	claims, ok := r.Context().Value(middlewares.ClaimsKey).(jwt.MapClaims)
	if !ok {
		return 0, http.ErrNoCookie // Or more appropriate error
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, http.ErrNoCookie
	}

	return int64(sub), nil
}
