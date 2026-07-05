package course

import (
	"backend/internal/middlewares"
	"backend/internal/modules/logs"
	"backend/internal/modules/users"
	"backend/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

func normalizeTime(t string) string {
	if len(t) == 5 {
		return t + ":00"
	}
	return t
}

func deriveStatus(statut string, actif bool) (string, bool) {
	if statut != "" {
		if statut == "brouillon" {
			return "brouillon", false
		}
		return "pending", false
	}
	if actif {
		return "approved", true
	}
	return "pending", false
}

type sessionInput struct {
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

func buildSessions(rawSessions []sessionInput, fallbackDate, fallbackStart, fallbackEnd string) ([]CourseSession, error) {
	if len(rawSessions) == 0 && fallbackDate != "" {
		rawSessions = []sessionInput{{Date: fallbackDate, StartTime: fallbackStart, EndTime: fallbackEnd}}
	}
	if len(rawSessions) == 0 {
		return nil, fmt.Errorf("au moins une date de session est requise")
	}

	sessions := make([]CourseSession, 0, len(rawSessions))
	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	nowMicros := int64(now.Hour())*3600*1_000_000 + int64(now.Minute())*60*1_000_000 + int64(now.Second())*1_000_000

	for _, raw := range rawSessions {
		var sess CourseSession
		if err := sess.SessionDate.Scan(raw.Date); err != nil || !sess.SessionDate.Valid {
			return nil, fmt.Errorf("date de session invalide")
		}
		if err := sess.StartTime.Scan(normalizeTime(raw.StartTime)); err != nil || !sess.StartTime.Valid {
			return nil, fmt.Errorf("heure de début de session invalide")
		}
		if err := sess.EndTime.Scan(normalizeTime(raw.EndTime)); err != nil || !sess.EndTime.Valid {
			return nil, fmt.Errorf("heure de fin de session invalide")
		}

		if sess.SessionDate.Time.Before(today) {
			return nil, fmt.Errorf("une date de session ne peut pas être dans le passé")
		}
		if sess.SessionDate.Time.Equal(today) && sess.StartTime.Microseconds < nowMicros {
			return nil, fmt.Errorf("une heure de session ne peut pas être dans le passé")
		}
		if sess.StartTime.Microseconds >= sess.EndTime.Microseconds {
			return nil, fmt.Errorf("l'heure de fin doit être strictement supérieure à l'heure de début pour chaque session")
		}

		sessions = append(sessions, sess)
	}

	sort.Slice(sessions, func(i, j int) bool {
		if !sessions[i].SessionDate.Time.Equal(sessions[j].SessionDate.Time) {
			return sessions[i].SessionDate.Time.Before(sessions[j].SessionDate.Time)
		}
		return sessions[i].StartTime.Microseconds < sessions[j].StartTime.Microseconds
	})

	return sessions, nil
}

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

func (h *Handler) GetAllForAdmin(w http.ResponseWriter, r *http.Request) {
	courses, err := h.service.GetAllForAdmin()
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
		Name        string         `json:"nom"`
		Title       string         `json:"titre"`
		Description string         `json:"description"`
		Prix        float64        `json:"prix"`
		Price       float64        `json:"price"`
		Categorie   string         `json:"categorie"`
		Actif       bool           `json:"actif"`
		Statut      string         `json:"statut"`
		Date        string         `json:"date"`
		StartTime   string         `json:"start_time"`
		EndTime     string         `json:"end_time"`
		MaxCapacity *int32         `json:"max_capacity"`
		Type        string         `json:"type"`
		SessionLink string         `json:"session_link"`
		Sessions    []sessionInput `json:"sessions"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := input.Name
	if name == "" {
		name = input.Title
	}
	status, approved := deriveStatus(input.Statut, input.Actif)
	price := input.Prix
	if price == 0 {
		price = input.Price
	}

	courseType := CourseType(input.Type)
	if courseType != EnLigne {
		courseType = Presentiel
	}

	sessions, err := buildSessions(input.Sessions, input.Date, input.StartTime, input.EndTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c := Course{
		Name:        name,
		Description: input.Description,
		Status:      status,
		Approved:    approved,
		Type:        courseType,
		SessionLink: pgtype.Text{String: input.SessionLink, Valid: input.SessionLink != ""},
		Date:        sessions[0].SessionDate,
		EndDate:     sessions[len(sessions)-1].SessionDate,
		StartTime:   sessions[0].StartTime,
		EndTime:     sessions[0].EndTime,
	}
	if input.MaxCapacity != nil {
		c.MaxCapacity = pgtype.Int4{Int32: *input.MaxCapacity, Valid: true}
	}
	c.Price.UnmarshalJSON([]byte(fmt.Sprintf("%f", price)))
	c.CreatedBy = pgtype.Int8{Int64: int64(sub), Valid: true}

	id, err := h.service.Create(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.service.ReplaceSessions(id, sessions); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logs.AddFromRequest(r, "Création de formation", fmt.Sprintf("Formation #%d: %s", id.Int64, name), "info")

	res := OffreFrontend{
		Id:          id.Int64,
		Nom:         name,
		Categorie:   input.Categorie,
		Description: input.Description,
		Actif:       approved,
		Prix:        price,
		Type:        string(courseType),
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
		Name        string         `json:"nom"`
		Title       string         `json:"titre"`
		Description string         `json:"description"`
		Prix        float64        `json:"prix"`
		Price       float64        `json:"price"`
		Categorie   string         `json:"categorie"`
		Actif       bool           `json:"actif"`
		Statut      string         `json:"statut"`
		Date        string         `json:"date"`
		StartTime   string         `json:"start_time"`
		EndTime     string         `json:"end_time"`
		MaxCapacity *int32         `json:"max_capacity"`
		Type        string         `json:"type"`
		SessionLink string         `json:"session_link"`
		Sessions    []sessionInput `json:"sessions"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := input.Name
	if name == "" {
		name = input.Title
	}
	status, approved := deriveStatus(input.Statut, input.Actif)
	price := input.Prix
	if price == 0 {
		price = input.Price
	}

	courseType := CourseType(input.Type)
	if courseType != EnLigne {
		courseType = Presentiel
	}

	sessions, err := buildSessions(input.Sessions, input.Date, input.StartTime, input.EndTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c := Course{
		Name:        name,
		Description: input.Description,
		Status:      status,
		Approved:    approved,
		Type:        courseType,
		SessionLink: pgtype.Text{String: input.SessionLink, Valid: input.SessionLink != ""},
		Date:        sessions[0].SessionDate,
		EndDate:     sessions[len(sessions)-1].SessionDate,
		StartTime:   sessions[0].StartTime,
		EndTime:     sessions[0].EndTime,
	}
	if input.MaxCapacity != nil {
		c.MaxCapacity = pgtype.Int4{Int32: *input.MaxCapacity, Valid: true}
	}
	c.Price.UnmarshalJSON([]byte(fmt.Sprintf("%f", price)))

	if err := h.service.Update(courseId, c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.service.ReplaceSessions(courseId, sessions); err != nil {
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
		Type:        string(updated.Type),
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

	if err := h.service.Approve(pgtype.Int8{Int64: id, Valid: true}, pgtype.Int8{Int64: int64(sub), Valid: true}); err != nil {
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

	courses, err := h.service.GetCoursesByCreator(pgtype.Int8{Int64: int64(sub), Valid: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func (h *Handler) UploadDocument(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	courseId := pgtype.Int8{Int64: id, Valid: true}

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
	role, _ := claims["role"].(string)

	existing, err := h.service.GetById(courseId)
	if err != nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}
	if role != "admin" && existing.CreatedBy.Int64 != int64(sub) {
		http.Error(w, "Forbidden: you do not own this course", http.StatusForbidden)
		return
	}

	filename, originalName, err := utils.SaveDocument(r, "document")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	docId, err := h.service.CreateDocument(CourseDocument{
		CourseId:     courseId,
		Filename:     filename,
		OriginalName: originalName,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logs.AddFromRequest(r, "Ajout document formation", fmt.Sprintf("Formation #%d: document %s ajouté", id, originalName), "info")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": docId.Int64, "filename": filename, "original_name": originalName})
}

func (h *Handler) GetDocuments(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	courseId := pgtype.Int8{Int64: id, Valid: true}

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
	role, _ := claims["role"].(string)

	existing, err := h.service.GetById(courseId)
	if err != nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}

	if role != "admin" && existing.CreatedBy.Int64 != int64(sub) {
		enrolled, err := h.service.IsUserEnrolled(courseId, pgtype.Int8{Int64: int64(sub), Valid: true})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !enrolled {
			http.Error(w, "Forbidden: you must join this course to view its documents", http.StatusForbidden)
			return
		}
	}

	docs, err := h.service.GetDocumentsByCourseId(courseId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(docs)
}

func (h *Handler) GetSessions(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	sessions, err := h.service.GetSessionsByCourseId(pgtype.Int8{Int64: id, Valid: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}

func (h *Handler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("docId")
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
	role, _ := claims["role"].(string)

	doc, err := h.service.GetDocumentById(pgtype.Int8{Int64: id, Valid: true})
	if err != nil {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	existingCourse, err := h.service.GetById(doc.CourseId)
	if err != nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}
	if role != "admin" && existingCourse.CreatedBy.Int64 != int64(sub) {
		http.Error(w, "Forbidden: you do not own this course", http.StatusForbidden)
		return
	}

	if err := h.service.DeleteDocument(pgtype.Int8{Int64: id, Valid: true}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
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
