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

func userIdFromContext(r *http.Request) pgtype.Int8 {
	if claims, ok := r.Context().Value(middlewares.ClaimsKey).(jwt.MapClaims); ok {
		if sub, ok := claims["sub"].(float64); ok {
			return pgtype.Int8{Int64: int64(sub), Valid: true}
		}
	}
	return pgtype.Int8{}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	posts, err := h.service.GetAll(r.URL.Query().Get("type"), userIdFromContext(r), r.URL.Query().Get("sort"))

	if err != nil {
		fmt.Println(err.Error())
	}

	res, _ := json.Marshal(posts)
	fmt.Fprintf(w, "%s", string(res))
}

func (h *Handler) GetAllActualites(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	posts, err := h.service.GetAll(string(Actualite), userIdFromContext(r), r.URL.Query().Get("sort"))

	if err != nil {
		fmt.Println(err.Error())
	}

	res, _ := json.Marshal(posts)
	fmt.Fprintf(w, "%s", string(res))
}

func (h *Handler) GetAllConseils(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	posts, err := h.service.GetAll(string(Conseil), userIdFromContext(r), r.URL.Query().Get("sort"))

	if err != nil {
		fmt.Println(err.Error())
	}

	res, _ := json.Marshal(posts)
	fmt.Fprintf(w, "%s", string(res))
}

func (h *Handler) Vote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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

	idStr := r.PathValue("id")
	idInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var input struct {
		Type string `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	updated, err := h.service.Vote(
		pgtype.Int8{Int64: idInt, Valid: true},
		pgtype.Int8{Int64: int64(sub), Valid: true},
		input.Type,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(updated)
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

func (h *Handler) createWithType(w http.ResponseWriter, r *http.Request, newsType NewsType) {
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
	newsDto.Type = newsType

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

func (h *Handler) CreateActualite(w http.ResponseWriter, r *http.Request) {
	h.createWithType(w, r, Actualite)
}

func (h *Handler) CreateConseil(w http.ResponseWriter, r *http.Request) {
	h.createWithType(w, r, Conseil)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")
	idInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(middlewares.ClaimsKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sub, ok := claims["sub"].(float64)
	role, _ := claims["role"].(string)
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
		return
	}

	newsId := pgtype.Int8{Int64: idInt, Valid: true}
	existing, err := h.service.GetById(newsId)
	if err != nil {
		http.Error(w, "News not found", http.StatusNotFound)
		return
	}

	if role != "admin" && existing.CreatedBy.Int64 != int64(sub) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	existing.Title = input.Title
	existing.Content = input.Content

	if err := h.service.Update(newsId, *existing); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(existing)
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
