package course

import (
	"backend/internal/middlewares"
	"backend/internal/modules/users"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type Handler struct {
	service     *Service
	userService *users.Service
}

func NewHandler(service *Service, userService *users.Service) *Handler {
	return &Handler{service: service, userService: userService}
}

func (h *Handler) GetAllApprovedCatalogue(w http.ResponseWriter, r *http.Request) {
	offres, err := h.service.GetAllApprovedCatalogue()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(offres)
}

func (h *Handler) GetAllCatalogue(w http.ResponseWriter, r *http.Request) {
	offres, err := h.service.GetAllCatalogue()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(offres)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	courses, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	course, err := h.service.GetById(pgtype.Int8{Int64: id, Valid: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(course)
}

func (h *Handler) GetUserCourses(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("user_id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	courses, err := h.service.GetUserCourses(pgtype.Int8{Int64: id, Valid: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
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

	var input struct {
		Name        string  `json:"nom"`
		Description string  `json:"description"`
		Prix        float64 `json:"prix"`
		Categorie   string  `json:"categorie"`
		Actif       bool    `json:"actif"`
		Date        string  `json:"date"`
		StartTime   string  `json:"start_time"`
		EndTime     string  `json:"end_time"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c := Course{
		Name:        input.Name,
		Description: input.Description,
		Approved:    input.Actif,
	}
	c.Price.UnmarshalJSON([]byte(fmt.Sprintf("%f", input.Prix)))
	c.CreatedBy = pgtype.Int8{Int64: int64(sub), Valid: true}
	c.Date.Scan(input.Date)
	c.StartTime.Scan(input.StartTime)
	c.EndTime.Scan(input.EndTime)

	// Validation: date must not be in the past
	if c.Date.Valid && c.Date.Time.Before(time.Now().Truncate(24*time.Hour)) {
		http.Error(w, "La date de la formation ne peut pas être dans le passé", http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	res := OffreFrontend{
		Id:          id.Int64,
		Nom:         input.Name,
		Categorie:   input.Categorie,
		Description: input.Description,
		Actif:       input.Actif,
		Prix:        input.Prix,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

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

	courseId := pgtype.Int8{Int64: id, Valid: true}
	existing, err := h.service.GetById(courseId)
	if err != nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}

	role, _ := claims["role"].(string)

	if role != "admin" && existing.CreatedBy.Int64 != int64(sub) {
		http.Error(w, "Forbidden: you do not own this course", http.StatusForbidden)
		return
	}

	var input struct {
		Name        string  `json:"nom"`
		Description string  `json:"description"`
		Prix        float64 `json:"prix"`
		Categorie   string  `json:"categorie"`
		Actif       bool    `json:"actif"`
		Date        string  `json:"date"`
		StartTime   string  `json:"start_time"`
		EndTime     string  `json:"end_time"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c := Course{
		Name:        input.Name,
		Description: input.Description,
		Approved:    input.Actif,
	}
	c.Price.UnmarshalJSON([]byte(fmt.Sprintf("%f", input.Prix)))
	c.Date.Scan(input.Date)
	c.StartTime.Scan(input.StartTime)
	c.EndTime.Scan(input.EndTime)

	// Validation: date must not be in the past
	if c.Date.Valid && c.Date.Time.Before(time.Now().Truncate(24*time.Hour)) {
		http.Error(w, "La date de la formation ne peut pas être dans le passé", http.StatusBadRequest)
		return
	}

	if err := h.service.Update(courseId, c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updated, err := h.service.GetById(courseId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := OffreFrontend{
		Id:          updated.Id.Int64,
		Nom:         updated.Name,
		Categorie:   input.Categorie,
		Description: updated.Description,
		Actif:       updated.Approved,
	}
	f, _ := updated.Price.Float64Value()
	res.Prix = f.Float64

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) Approve(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.service.Approve(pgtype.Int8{Int64: id, Valid: true}, pgtype.Int8{Int64: 1, Valid: true}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Disapprove(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.service.Disapprove(pgtype.Int8{Int64: id, Valid: true}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

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

	courseId := pgtype.Int8{Int64: id, Valid: true}
	existing, err := h.service.GetById(courseId)
	if err != nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}

	role, _ := claims["role"].(string)

	if role != "admin" && existing.CreatedBy.Int64 != int64(sub) {
		http.Error(w, "Forbidden: you do not own this course", http.StatusForbidden)
		return
	}

	if err := h.service.Delete(courseId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
