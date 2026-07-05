package advertisement

import (
	"backend/internal/middlewares"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
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

func (h *Handler) GetMyAds(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	ads, err := h.service.GetByAnnouncerId(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ads)
}

func (h *Handler) GetAllPubs(w http.ResponseWriter, r *http.Request) {
	pubs, err := h.service.GetAllPubs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pubs)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	ads, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ads)
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	ad, err := h.service.GetById(pgtype.Int8{Int64: id, Valid: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ad)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var ad Advertisement
	if err := json.NewDecoder(r.Body).Decode(&ad); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ad.StartDate.Valid && ad.StartDate.Time.Before(time.Now().Truncate(24*time.Hour)) {
		http.Error(w, "La date de début de la publicité ne peut pas être dans le passé", http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(ad)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ad.Id = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ad)
}

func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	adminId, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var body struct {
		Statut string `json:"statut"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dbStatus := "pending"
	if body.Statut == "active" {
		dbStatus = "validated"
	}
	if err := h.service.UpdateStatus(pgtype.Int8{Int64: id, Valid: true}, dbStatus, pgtype.Int8{Int64: adminId, Valid: true}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Approve(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	adminId, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.service.UpdateStatus(pgtype.Int8{Int64: id, Valid: true}, "validated", pgtype.Int8{Int64: adminId, Valid: true}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Reject(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	adminId, err := getUserIdFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.service.UpdateStatus(pgtype.Int8{Int64: id, Valid: true}, "rejected", pgtype.Int8{Int64: adminId, Valid: true}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.service.Delete(pgtype.Int8{Int64: id, Valid: true}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
