package post

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

	posts, err := h.service.GetAll()

	if err != nil {
		fmt.Println(err.Error())
	}

	res, _ := json.Marshal(posts)
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

	post, err := h.service.GetById(pgtype.Int8{Int64: idInt, Valid: true})

	res, _ := json.Marshal(post)
	fmt.Fprintf(w, "%s", string(res))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var postDto Post

	err := json.NewDecoder(r.Body).Decode(&postDto)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(postDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	postDto.Id = id
	res, _ := json.Marshal(postDto)
	fmt.Fprintf(w, "%s", string(res))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")

	idInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req struct {
		Content string `json:"content"`
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateContent(pgtype.Int8{Int64: idInt, Valid: true}, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "post updated successfully"}`)
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
	fmt.Fprintf(w, `{"message": "post deleted successfully"}`)
}
