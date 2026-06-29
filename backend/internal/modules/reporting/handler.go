package reporting

import (
	"backend/internal/middlewares"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ExportAuditPDF(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	if startStr == "" || endStr == "" {
		http.Error(w, "Les parametres 'start' et 'end' (format YYYY-MM-DD) sont obligatoires", http.StatusBadRequest)
		return
	}

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		http.Error(w, "Format de date 'start' invalide (attendu: YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		http.Error(w, "Format de date 'end' invalide (attendu: YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	if start.After(end) {
		http.Error(w, "La date de debut doit être anterieure a la date de fin", http.StatusBadRequest)
		return
	}

	end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	pdfBytes, err := h.service.GenerateAuditPDF(r.Context(), start, end)
	if err != nil {
		http.Error(w, "Erreur lors de la generation du PDF: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=audit_depots.pdf")
	w.Write(pdfBytes)
}

func (h *Handler) GetActorStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetActorStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

func (h *Handler) GetPrestationStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetPrestationStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

func (h *Handler) GetUserPredictions(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = 10
	}

	predictions, err := h.service.GetUserPredictions(r.Context(), page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(predictions)
}

func (h *Handler) GetMLStatus(w http.ResponseWriter, r *http.Request) {
	status, err := h.service.GetMLStatus(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(status)
}

func (h *Handler) GetPredictionDistribution(w http.ResponseWriter, r *http.Request) {
	dist, err := h.service.repo.GetPredictionDistribution(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(dist)
}

func (h *Handler) GetSalarieStats(w http.ResponseWriter, r *http.Request) {
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

	stats, err := h.service.GetSalarieStats(r.Context(), int64(sub))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
