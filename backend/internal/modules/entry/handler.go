package entry

import (
	"github.com/jackc/pgx/v5/pgtype"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	items, err := h.service.GetAll()
	if err != nil {
		fmt.Println(err.Error())
	}

	res, _ := json.Marshal(items)
	fmt.Fprintf(w, "%s", string(res))
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")
	idInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item, err := h.service.GetById(pgtype.Int8{Int64: idInt, Valid: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	res, _ := json.Marshal(item)
	fmt.Fprintf(w, "%s", string(res))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var dto Entry
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dto.Id = id
	res, _ := json.Marshal(dto)
	fmt.Fprintf(w, "%s", string(res))
}

func (h *Handler) DeleteById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")
	idInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.Delete(pgtype.Int8{Int64: idInt, Valid: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "entry deleted successfully"}`)
}
