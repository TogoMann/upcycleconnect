package planning

import (
	"backend/internal/middlewares"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

func clockToPgTime(s string) pgtype.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return pgtype.Time{Valid: false}
	}
	
	s = strings.ReplaceAll(s, "h", ":")
	s = strings.ReplaceAll(s, "H", ":")
	
	isPM := false
	isAM := false
	if strings.HasSuffix(strings.ToUpper(s), "PM") {
		isPM = true
		s = strings.TrimSpace(s[:len(s)-2])
	} else if strings.HasSuffix(strings.ToUpper(s), "AM") {
		isAM = true
		s = strings.TrimSpace(s[:len(s)-2])
	}
	
	parts := strings.Split(s, ":")
	if len(parts) < 2 {
		return pgtype.Time{Valid: false}
	}
	
	hStr := strings.TrimSpace(parts[0])
	mStr := strings.TrimSpace(parts[1])
	
	h, errH := strconv.Atoi(hStr)
	m, errM := strconv.Atoi(mStr)
	if errH != nil || errM != nil {
		return pgtype.Time{Valid: false}
	}
	
	sec := 0
	if len(parts) >= 3 {
		secStr := strings.TrimSpace(parts[2])
		if idx := strings.IndexAny(secStr, ".,-+"); idx != -1 {
			secStr = secStr[:idx]
		}
		if v, err := strconv.Atoi(secStr); err == nil {
			sec = v
		}
	}
	
	if isPM && h < 12 {
		h += 12
	} else if isAM && h == 12 {
		h = 0
	}
	
	if h < 0 || h > 23 || m < 0 || m > 59 || sec < 0 || sec > 59 {
		return pgtype.Time{Valid: false}
	}
	
	micros := (int64(h)*3600 + int64(m)*60 + int64(sec)) * 1_000_000
	return pgtype.Time{Microseconds: micros, Valid: true}
}

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

func (h *Handler) GetAllPlannings(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetAllPlannings()
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

	now := time.Now()
	if loc, err := time.LoadLocation("Europe/Paris"); err == nil {
		now = now.In(loc)
	}
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if e.Date.Time.Before(today) {
		http.Error(w, "La date de l'événement ne peut pas être dans le passé", http.StatusBadRequest)
		return
	}

	e.StartTime = clockToPgTime(input.StartTime)
	e.EndTime = clockToPgTime(input.EndTime)

	if !e.StartTime.Valid {
		e.StartTime = clockToPgTime("00:00:00")
	}
	if !e.EndTime.Valid {
		e.EndTime = clockToPgTime("23:59:59")
	}

	if e.Date.Valid && e.StartTime.Valid {
		if e.Date.Time.Equal(today) {
			nowMicros := int64(now.Hour())*3600*1_000_000 + int64(now.Minute())*60*1_000_000 + int64(now.Second())*1_000_000
			if e.StartTime.Microseconds < nowMicros {
				http.Error(w, "L'heure de début ne peut pas être dans le passé", http.StatusBadRequest)
				return
			}
		}
	}

	if e.StartTime.Valid && e.EndTime.Valid && e.StartTime.Microseconds >= e.EndTime.Microseconds {
		http.Error(w, "L'heure de fin doit être strictement supérieure à l'heure de début", http.StatusBadRequest)
		return
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
