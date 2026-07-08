package entry

import (
	"backend/internal/middlewares"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

func (h *Handler) GetAllDepots(w http.ResponseWriter, r *http.Request) {
	depots, err := h.service.GetAllDepots()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("%#v\n", depots)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(depots)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	entries, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	entry, err := h.service.GetById(pgtype.Int8{Int64: id, Valid: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entry)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
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

	var input struct {
		Schedule string `json:"schedule"`
		Start    string `json:"start"`
		Ending   string `json:"ending"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var e Entry
	e.CreatedBy = pgtype.Int8{Int64: int64(sub), Valid: true}
	e.Schedule.Scan(input.Schedule)
	e.Start = clockToPgTime(input.Start)
	if input.Ending != "" {
		e.Ending = clockToPgTime(input.Ending)
	}

	now := time.Now()
	if loc, err := time.LoadLocation("Europe/Paris"); err == nil {
		now = now.In(loc)
	}
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if e.Schedule.Valid && e.Schedule.Time.Before(today) {
		http.Error(w, "La date de dépôt ne peut pas être dans le passé", http.StatusBadRequest)
		return
	}

	if e.Schedule.Valid && e.Start.Valid {
		if e.Schedule.Time.Equal(today) {
			nowMicros := int64(now.Hour())*3600*1_000_000 + int64(now.Minute())*60*1_000_000 + int64(now.Second())*1_000_000
			if e.Start.Microseconds < nowMicros {
				http.Error(w, "L'heure de début ne peut pas être dans le passé", http.StatusBadRequest)
				return
			}
		}
	}

	id, err := h.service.Create(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	e.Id = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(e)
}

func (h *Handler) DeleteById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.service.Delete(pgtype.Int8{Int64: id, Valid: true}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ValiderDepot(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.service.ValidateDepot(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) EnvoyerCode(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "Code envoyé"}`)
}
