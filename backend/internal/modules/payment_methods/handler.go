package payment_methods

import (
	"backend/internal/middlewares"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetMyPaymentMethods(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	methods, err := h.service.GetByUserId(pgtype.Int8{Int64: userId, Valid: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(methods)
}

func (h *Handler) HasPaymentMethod(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	has, err := h.service.HasPaymentMethod(pgtype.Int8{Int64: userId, Valid: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"has_payment_method": has})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var pm PaymentMethod
	if err := json.NewDecoder(r.Body).Decode(&pm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pm.UserId = pgtype.Int8{Int64: userId, Valid: true}
	id, err := h.service.Create(pm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pm.Id = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pm)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.service.DeleteByIdAndUserId(pgtype.Int8{Int64: id, Valid: true}, pgtype.Int8{Int64: userId, Valid: true}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) SetDefault(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.service.SetDefault(pgtype.Int8{Int64: id, Valid: true}, pgtype.Int8{Int64: userId, Valid: true}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	methods, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(methods)
}

func getUserIdFromContext(r *http.Request) (int64, error) {
	claims, ok := r.Context().Value(middlewares.ClaimsKey).(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("no claims found")
	}
	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid sub in claims")
	}
	return int64(sub), nil
}
