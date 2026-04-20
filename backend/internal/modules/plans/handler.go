package plans

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	plans, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	frontendPlans := make([]PlanFrontend, len(plans))
	for i, p := range plans {
		f, _ := p.Price.Float64Value()
		frontendPlans[i] = PlanFrontend{
			Id:           p.Id.Int64,
			Name:         p.Name,
			Description:  p.Description,
			Price:        f.Float64,
			BillingCycle: p.BillingCycle,
			Features:     p.Features,
			IsActive:     p.IsActive,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(frontendPlans)
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	p, err := h.service.GetById(pgtype.Int8{Int64: id, Valid: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	f, _ := p.Price.Float64Value()
	res := PlanFrontend{
		Id:           p.Id.Int64,
		Name:         p.Name,
		Description:  p.Description,
		Price:        f.Float64,
		BillingCycle: p.BillingCycle,
		Features:     p.Features,
		IsActive:     p.IsActive,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input PlanFrontend
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p := Plan{
		Name:         input.Name,
		Description:  input.Description,
		BillingCycle: input.BillingCycle,
		Features:     input.Features,
		IsActive:     input.IsActive,
	}
	p.Price.UnmarshalJSON([]byte(fmt.Sprintf("%f", input.Price)))

	id, err := h.service.Create(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	input.Id = id.Int64
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(input)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	var input PlanFrontend
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p := Plan{
		Name:         input.Name,
		Description:  input.Description,
		BillingCycle: input.BillingCycle,
		Features:     input.Features,
		IsActive:     input.IsActive,
	}
	p.Price.UnmarshalJSON([]byte(fmt.Sprintf("%f", input.Price)))

	if err := h.service.Update(pgtype.Int8{Int64: id, Valid: true}, p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	input.Id = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(input)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.service.Delete(pgtype.Int8{Int64: id, Valid: true}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
