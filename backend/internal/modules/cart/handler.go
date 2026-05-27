package cart

import (
	"backend/internal/middlewares"
	"encoding/json"
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

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.ClaimsKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	userId := pgtype.Int8{Int64: int64(sub), Valid: true}
	items, err := h.service.GetByUserId(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.ClaimsKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	var body struct {
		ListingId int64 `json:"listing_id,omitempty"`
		EventId   int64 `json:"event_id,omitempty"`
		CourseId  int64 `json:"course_id,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	userId := pgtype.Int8{Int64: int64(sub), Valid: true}
	listingId := pgtype.Int8{Int64: body.ListingId, Valid: body.ListingId != 0}
	eventId := pgtype.Int8{Int64: body.EventId, Valid: body.EventId != 0}
	courseId := pgtype.Int8{Int64: body.CourseId, Valid: body.CourseId != 0}

	if err := h.service.Add(userId, listingId, eventId, courseId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Remove(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.ClaimsKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	itemType := r.PathValue("type")
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	userId := pgtype.Int8{Int64: int64(sub), Valid: true}
	var listingId, eventId, courseId pgtype.Int8

	switch itemType {
	case "listing":
		listingId = pgtype.Int8{Int64: id, Valid: true}
	case "event":
		eventId = pgtype.Int8{Int64: id, Valid: true}
	case "course":
		courseId = pgtype.Int8{Int64: id, Valid: true}
	default:
		http.Error(w, "Invalid item type", http.StatusBadRequest)
		return
	}

	if err := h.service.Remove(userId, listingId, eventId, courseId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Checkout(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.ClaimsKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	userId := pgtype.Int8{Int64: int64(sub), Valid: true}
	if err := h.service.Checkout(userId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Checkout successful"})
}
