package users

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

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims, ok := r.Context().Value(middlewares.ClaimsKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	username, ok := claims["username"].(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.service.GetByUsername(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, _ := json.Marshal(user)
	fmt.Fprintf(w, "%s", string(res))
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := h.service.GetAll()

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, _ := json.Marshal(users)
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

	user, err := h.service.GetById(pgtype.Int8{Int64: idInt, Valid: true})

	res, _ := json.Marshal(user)
	fmt.Fprintf(w, "%s", string(res))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var userDto User

	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(userDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	userDto.Id = id
	res, _ := json.Marshal(userDto)
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

	userId := pgtype.Int8{Int64: idInt, Valid: true}

	
	claims, ok := r.Context().Value(middlewares.ClaimsKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	role, _ := claims["role"].(string)
	username, _ := claims["username"].(string)

	existingUser, err := h.service.GetById(userId)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if role != string(Admin) && existingUser.Username != username {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var updateData map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	
	if val, ok := updateData["username"].(string); ok {
		existingUser.Username = val
	}
	if val, ok := updateData["first_name"].(string); ok {
		existingUser.FirstName = val
	}
	if val, ok := updateData["last_name"].(string); ok {
		existingUser.LastName = val
	}
	if val, ok := updateData["email"].(string); ok {
		existingUser.Email = val
	}
	if val, ok := updateData["role"].(string); ok {
		existingUser.Role = UserRole(val)
	}
	if val, ok := updateData["language_preference"].(string); ok {
		existingUser.LanguagePreference = val
	}

	
	if role != string(Admin) {
		
		
		
		
		
		
		
	}
	
	
	if role != string(Admin) {
		
		originalUser, _ := h.service.GetById(userId)
		existingUser.Role = originalUser.Role
		existingUser.Username = originalUser.Username
	}

	err = h.service.Update(userId, *existingUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "user updated successfully"}`)
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
	fmt.Fprintf(w, `{"message": "user deleted successfully"}`)
}

func (h *Handler) GetScore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")
	idInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	score, err := h.service.GetScore(pgtype.Int8{Int64: idInt, Valid: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, `{"score": %d}`, score)
}

func (h *Handler) UpdateTutorialSeen(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")
	idInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.UpdateTutorialSeen(pgtype.Int8{Int64: idInt, Valid: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, `{"message": "tutorial status updated"}`)
}
