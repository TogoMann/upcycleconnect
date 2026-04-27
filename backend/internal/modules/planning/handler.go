package planning

import (
	"backend/internal/middlewares"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetMyPlanning(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.ClaimsKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	items, err := h.service.GetUserPlanning(pgtype.Int8{Int64: int64(sub), Valid: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *Handler) CreatePersonalEvent(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.ClaimsKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Date        string `json:"date"`
		StartTime   string `json:"start_time"`
		EndTime     string `json:"end_time"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	e := PersonalEvent{
		UserId:      pgtype.Int8{Int64: int64(sub), Valid: true},
		Title:       input.Title,
		Description: input.Description,
	}
	e.Date.Scan(input.Date)
	
	// Validation: date must not be in the past
	if e.Date.Time.Before(time.Now().Truncate(24 * time.Hour)) {
		http.Error(w, "La date de l'événement ne peut pas être dans le passé", http.StatusBadRequest)
		return
	}

	e.StartTime.Scan(input.StartTime)
	e.EndTime.Scan(input.EndTime)

	if !e.StartTime.Valid {
		e.StartTime.Scan("00:00:00")
	}
	if !e.EndTime.Valid {
		e.EndTime.Scan("23:59:59")
	}

	id, err := h.service.CreatePersonalEvent(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	e.Id = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(e)
}

func (h *Handler) DeletePersonalEvent(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.ClaimsKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	err := h.service.DeletePersonalEvent(pgtype.Int8{Int64: id, Valid: true}, pgtype.Int8{Int64: int64(sub), Valid: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
