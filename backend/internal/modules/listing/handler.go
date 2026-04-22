package listing

import (
	"backend/internal/middlewares"
	"backend/internal/modules/users"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

func (h *Handler) Approve(w http.ResponseWriter, r *http.Request) {
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

	username, _ := claims["username"].(string)
	admin, err := h.userService.GetByUsername(username)
	if err != nil {
		http.Error(w, "Admin user not found", http.StatusInternalServerError)
		return
	}

	err = h.service.Approve(pgtype.Int8{Int64: idInt, Valid: true}, admin.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, `{"message": "listing approved successfully"}`)
}

func (h *Handler) Disapprove(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")
	idInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.Disapprove(pgtype.Int8{Int64: idInt, Valid: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, `{"message": "listing disapproved successfully"}`)
}

func (h *Handler) GetAllApproved(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	listings, err := h.service.GetAllApproved()
	if err != nil {
		fmt.Println(err.Error())
	}
	res, _ := json.Marshal(listings)
	fmt.Fprintf(w, "%s", string(res))
}

func (h *Handler) GetByUserId(w http.ResponseWriter, r *http.Request) {
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

	listings, err := h.service.GetByUserId(pgtype.Int8{Int64: int64(sub), Valid: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(listings)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	listings, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if listings == nil {
		listings = []Listing{}
	}
	json.NewEncoder(w).Encode(listings)
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")

	idInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	listing, err := h.service.GetById(pgtype.Int8{Int64: idInt, Valid: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(listing)
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
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Category    string  `json:"category"`
		CityId      int64   `json:"city_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	listingDto := Listing{
		Name:        input.Name,
		Description: input.Description,
		Category:    ListingCategory(input.Category),
		CityId:      pgtype.Int8{Int64: input.CityId, Valid: input.CityId > 0},
		CreatedBy:   pgtype.Int8{Int64: int64(sub), Valid: true},
		Status:      Active,
		Approved:    false,
	}
	listingDto.Price.UnmarshalJSON([]byte(fmt.Sprintf("%f", input.Price)))

	id, err := h.service.Create(listingDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	listingDto.Id = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(listingDto)
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
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
		return
	}

	existing, err := h.service.GetById(pgtype.Int8{Int64: idInt, Valid: true})
	if err != nil {
		http.Error(w, "Listing not found", http.StatusNotFound)
		return
	}

	role, _ := claims["role"].(string)

	if role != "admin" && existing.CreatedBy.Int64 != int64(sub) {
		http.Error(w, "Forbidden: you do not own this listing", http.StatusForbidden)
		return
	}

	var input struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Category    string  `json:"category"`
		Status      string  `json:"status"`
		CityId      int64   `json:"city_id"`
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	dto := Listing{
		Name:        input.Name,
		Description: input.Description,
		Category:    ListingCategory(input.Category),
		Status:      ListingStatus(input.Status),
		CityId:      pgtype.Int8{Int64: input.CityId, Valid: input.CityId > 0},
	}
	dto.Price.UnmarshalJSON([]byte(fmt.Sprintf("%f", input.Price)))

	err = h.service.Update(pgtype.Int8{Int64: idInt, Valid: true}, dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "listing updated successfully"}`)
}

func (h *Handler) DeleteById(w http.ResponseWriter, r *http.Request) {
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
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
		return
	}

	existing, err := h.service.GetById(pgtype.Int8{Int64: idInt, Valid: true})
	if err != nil {
		http.Error(w, "Listing not found", http.StatusNotFound)
		return
	}

	role, _ := claims["role"].(string)

	if role != "admin" && existing.CreatedBy.Int64 != int64(sub) {
		http.Error(w, "Forbidden: you do not own this listing", http.StatusForbidden)
		return
	}

	err = h.service.Delete(pgtype.Int8{Int64: idInt, Valid: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "listing deleted successfully"}`)
}
