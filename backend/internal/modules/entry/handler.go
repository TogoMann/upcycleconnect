package entry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

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
	var e Entry
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if e.Schedule.Valid && e.Schedule.Time.Before(time.Now().Truncate(24*time.Hour)) {
		http.Error(w, "La date de dépôt ne peut pas être dans le passé", http.StatusBadRequest)
		return
	}

	if e.Schedule.Valid && e.Start.Valid {
		now := time.Now()
		today := now.Truncate(24 * time.Hour)
		if e.Schedule.Time.Equal(today) {
			nowMicros := int64(now.Hour())*3600*1_000_000 + int64(now.Minute())*60*1_000_000 + int64(now.Second())*1_000_000
			if e.Start.Microseconds < nowMicros {
				http.Error(w, "L'heure de début ne peut pas être dans le passé", http.StatusBadRequest)
				return
			}
		}
	}

	if e.Start.Valid && e.Ending.Valid && e.Start.Microseconds >= e.Ending.Microseconds {
		http.Error(w, "L'heure de fin doit être strictement supérieure à l'heure de début", http.StatusBadRequest)
		return
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
