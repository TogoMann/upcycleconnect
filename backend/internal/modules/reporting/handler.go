package reporting

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
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

	// Ensure start is before end
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
