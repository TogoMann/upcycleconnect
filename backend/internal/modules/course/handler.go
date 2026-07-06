package course

import (
	"backend/internal/middlewares"
	"backend/internal/modules/logs"
	"backend/internal/modules/users"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Handler struct {
	service     *Service
	userService *users.Service
}

func NewHandler(service *Service, userService *users.Service) *Handler {
	return &Handler{service: service, userService: userService}
}

func normalizeTime(s string) string {
	if len(s) == 5 {
		return s + ":00"
	}
	return s
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
		Title       string  `json:"titre"`
		Description string  `json:"description"`
		Prix        float64 `json:"prix"`
		Price       float64 `json:"price"`
		Categorie   string  `json:"categorie"`
		Actif       bool    `json:"actif"`
		Statut      string  `json:"statut"`
		Date        string  `json:"date"`
		StartTime   string  `json:"start_time"`
		EndTime     string  `json:"end_time"`
		MaxCapacity *int32  `json:"max_capacity"`
		Duree       string  `json:"duree"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := input.Name
	if name == "" {
		name = input.Title
	}
	approved := input.Actif
	if input.Statut == "publiee" {
		approved = true
	} else if input.Statut == "brouillon" {
		approved = false
	}
	price := input.Prix
	if price == 0 {
		price = input.Price
	}
	if price < 0 {
		http.Error(w, "Le prix ne peut pas être négatif", http.StatusBadRequest)
		return
	}
	if input.MaxCapacity != nil && *input.MaxCapacity < 0 {
		http.Error(w, "La capacité maximale ne peut pas être négative", http.StatusBadRequest)
		return
	}
	desc := input.Description
	{

	}

	c := Course{
		Name:        name,
		Description: desc,
		Approved:    approved,
	}
	if input.MaxCapacity != nil {
		c.MaxCapacity = pgtype.Int4{Int32: *input.MaxCapacity, Valid: true}
	}
	c.Price.UnmarshalJSON([]byte(fmt.Sprintf("%f", price)))
	c.CreatedBy = pgtype.Int8{Int64: int64(sub), Valid: true}
	c.Date.Scan(input.Date)
	c.StartTime.Scan(normalizeTime(input.StartTime))
	c.EndTime.Scan(normalizeTime(input.EndTime))

	if c.Date.Valid && c.Date.Time.Before(time.Now().Truncate(24*time.Hour)) {
		http.Error(w, "La date de la formation ne peut pas être dans le passé", http.StatusBadRequest)
		return
	}

	if c.Date.Valid && c.StartTime.Valid {
		now := time.Now()
		today := now.Truncate(24 * time.Hour)
		if c.Date.Time.Equal(today) {
			nowMicros := int64(now.Hour())*3600*1_000_000 + int64(now.Minute())*60*1_000_000 + int64(now.Second())*1_000_000
			if c.StartTime.Microseconds < nowMicros {
				http.Error(w, "L'heure de début ne peut pas être dans le passé", http.StatusBadRequest)
				return
			}
		}
	}

	if c.StartTime.Valid && c.EndTime.Valid && c.StartTime.Microseconds >= c.EndTime.Microseconds {
		http.Error(w, "L'heure de fin doit être strictement supérieure à l'heure de début", http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(c, input.Categorie, input.Duree)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logs.AddFromRequest(r, "Création de formation", fmt.Sprintf("Formation #%d: %s", id.Int64, name), "info")

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
		Title       string  `json:"titre"`
		Description string  `json:"description"`
		Prix        float64 `json:"prix"`
		Price       float64 `json:"price"`
		Categorie   string  `json:"categorie"`
		Actif       bool    `json:"actif"`
		Statut      string  `json:"statut"`
		Date        string  `json:"date"`
		StartTime   string  `json:"start_time"`
		EndTime     string  `json:"end_time"`
		MaxCapacity *int32  `json:"max_capacity"`
		Duree       string  `json:"duree"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := input.Name
	if name == "" {
		name = input.Title
	}
	approved := input.Actif
	if input.Statut == "publiee" {
		approved = true
	} else if input.Statut == "brouillon" {
		approved = false
	}
	price := input.Prix
	if price == 0 {
		price = input.Price
	}
	if price < 0 {
		http.Error(w, "Le prix ne peut pas être négatif", http.StatusBadRequest)
		return
	}
	if input.MaxCapacity != nil && *input.MaxCapacity < 0 {
		http.Error(w, "La capacité maximale ne peut pas être négative", http.StatusBadRequest)
		return
	}
	desc := input.Description
	if input.Duree != "" {
		desc = desc + "\n\nDurée: " + input.Duree
	}

	c := Course{
		Name:        name,
		Description: desc,
		Approved:    approved,
	}
	if input.MaxCapacity != nil {
		c.MaxCapacity = pgtype.Int4{Int32: *input.MaxCapacity, Valid: true}
	}
	c.Price.UnmarshalJSON([]byte(fmt.Sprintf("%f", price)))
	c.Date.Scan(input.Date)
	c.StartTime.Scan(normalizeTime(input.StartTime))
	c.EndTime.Scan(normalizeTime(input.EndTime))

	if c.Date.Valid && c.Date.Time.Before(time.Now().Truncate(24*time.Hour)) {
		http.Error(w, "La date de la formation ne peut pas être dans le passé", http.StatusBadRequest)
		return
	}

	if c.Date.Valid && c.StartTime.Valid {
		now := time.Now()
		today := now.Truncate(24 * time.Hour)
		if c.Date.Time.Equal(today) {
			nowMicros := int64(now.Hour())*3600*1_000_000 + int64(now.Minute())*60*1_000_000 + int64(now.Second())*1_000_000
			if c.StartTime.Microseconds < nowMicros {
				http.Error(w, "L'heure de début ne peut pas être dans le passé", http.StatusBadRequest)
				return
			}
		}
	}

	if c.StartTime.Valid && c.EndTime.Valid && c.StartTime.Microseconds >= c.EndTime.Microseconds {
		http.Error(w, "L'heure de fin doit être strictement supérieure à l'heure de début", http.StatusBadRequest)
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

func (h *Handler) GetMyCourses(w http.ResponseWriter, r *http.Request) {
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

	courses, err := h.service.GetCreatedByUser(pgtype.Int8{Int64: int64(sub), Valid: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func (h *Handler) ProposeModification(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	var input struct {
		Comment     string  `json:"comment"`
		Name        string  `json:"nom"`
		Description string  `json:"description"`
		Prix        float64 `json:"prix"`
		Date        string  `json:"date"`
		StartTime   string  `json:"start_time"`
		EndTime     string  `json:"end_time"`
		MaxCapacity *int32  `json:"max_capacity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	courseId := pgtype.Int8{Int64: id, Valid: true}
	existing, err := h.service.GetById(courseId)
	if err != nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}

	if input.Name != "" {
		existing.Name = input.Name
	}
	if input.Description != "" {
		existing.Description = input.Description
	}
	if input.Prix != 0 {
		existing.Price.UnmarshalJSON([]byte(fmt.Sprintf("%f", input.Prix)))
	}
	if input.Date != "" {
		existing.Date.Scan(input.Date)
	}
	if input.StartTime != "" {
		existing.StartTime.Scan(normalizeTime(input.StartTime))
	}
	if input.EndTime != "" {
		existing.EndTime.Scan(normalizeTime(input.EndTime))
	}
	if input.MaxCapacity != nil {
		existing.MaxCapacity = pgtype.Int4{Int32: *input.MaxCapacity, Valid: true}
	}

	existing.Status = "needs_modification"
	existing.CorrectionComment = pgtype.Text{String: input.Comment, Valid: true}

	if err := h.service.Update(courseId, *existing); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "Modification proposal sent successfully"}`)
}
