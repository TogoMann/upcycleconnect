package news

import (
	"backend/internal/middlewares"
	"encoding/json"
	"fmt"
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
	claims, ok := r.Context().Value(middlewares.ClaimsKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
		return
	}

	var newsDto News
	err := json.NewDecoder(r.Body).Decode(&newsDto)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	newsDto.CreatedBy = pgtype.Int8{Int64: int64(sub), Valid: true}

	id, err := h.service.Create(newsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newsDto.Id = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newsDto)
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
	fmt.Fprintf(w, `{"message": "news deleted successfully"}`)
}
