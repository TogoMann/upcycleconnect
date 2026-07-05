package listing

import (
	"backend/internal/middlewares"
	"backend/internal/modules/logs"
	"backend/internal/modules/users"
	"backend/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *Handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	filename, err := utils.SaveImage(r, "image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"url": "/uploads/" + filename,
	})
}

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

	logs.AddFromRequest(r, "Approbation d'annonce", fmt.Sprintf("Annonce #%d approuvée", idInt), "info")

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

	logs.AddFromRequest(r, "Désapprobation d'annonce", fmt.Sprintf("Annonce #%d désapprouvée", idInt), "info")

	fmt.Fprintf(w, `{"message": "listing disapproved successfully"}`)
}

func (h *Handler) GetAllApproved(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	excludeUserId := pgtype.Int8{}
	if claims, ok := r.Context().Value(middlewares.ClaimsKey).(jwt.MapClaims); ok {
		if sub, ok := claims["sub"].(float64); ok {
			excludeUserId = pgtype.Int8{Int64: int64(sub), Valid: true}
		}
	}

	minLevel, _ := strconv.ParseInt(r.URL.Query().Get("min_level"), 10, 32)

	listings, err := h.service.GetAllApproved(excludeUserId, int32(minLevel))
	if err != nil {
		fmt.Println(err.Error())
	}
	for i := range listings {
		listings[i].SellerLevel = utils.LevelFromScore(listings[i].SellerScore)
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

	for i := range listings {
		listings[i].SellerLevel = utils.LevelFromScore(listings[i].SellerScore)
	}

	json.NewEncoder(w).Encode(listings)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = 10
	}

	listings, err := h.service.GetAll(page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

	listing.SellerLevel = utils.LevelFromScore(listing.SellerScore)

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
		Name          string  `json:"name"`
		Description   string  `json:"description"`
		Price         float64 `json:"price"`
		Category      string  `json:"category"`
		CityId        int64   `json:"city_id"`
		ImageUrl      string  `json:"image_url"`
		HandoffMode   string  `json:"handoff_mode"`
		Address       string  `json:"address"`
		Weight        float64 `json:"weight"`
		LockerId      int64   `json:"locker_id"`
		PhysicalState string  `json:"physical_state"`
		Size          string  `json:"size"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	handoffMode := ListingHandoffMode(input.HandoffMode)
	if handoffMode != Locker {
		handoffMode = HandDelivery
	}

	listingDto := Listing{
		Name:        input.Name,
		Description: input.Description,
		Category:    ListingCategory(input.Category),
		CityId:      pgtype.Int8{Int64: input.CityId, Valid: input.CityId > 0},
		CreatedBy:   pgtype.Int8{Int64: int64(sub), Valid: true},
		Status:      Active,
		Approved:    false,
		ImageUrl:    pgtype.Text{String: input.ImageUrl, Valid: input.ImageUrl != ""},
		HandoffMode: handoffMode,
		Address:     pgtype.Text{String: input.Address, Valid: input.Address != ""},
	}
	listingDto.Price.UnmarshalJSON([]byte(fmt.Sprintf("%f", input.Price)))
	listingDto.Weight.UnmarshalJSON([]byte(fmt.Sprintf("%f", input.Weight)))

	lockerId := pgtype.Int8{Int64: input.LockerId, Valid: input.LockerId > 0}
	id, err := h.service.Create(listingDto, lockerId, input.PhysicalState, input.Size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	listingDto.Id = id
	logs.AddFromRequest(r, "Création d'annonce", fmt.Sprintf("Annonce #%d: %s", id, listingDto.Name), "info")
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
		ImageUrl    string  `json:"image_url"`
		HandoffMode string  `json:"handoff_mode"`
		Address     string  `json:"address"`
		Weight      float64 `json:"weight"`
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	handoffMode := ListingHandoffMode(input.HandoffMode)
	if handoffMode != Locker {
		handoffMode = HandDelivery
	}

	dto := Listing{
		Name:        input.Name,
		Description: input.Description,
		Category:    ListingCategory(input.Category),
		Status:      ListingStatus(input.Status),
		CityId:      pgtype.Int8{Int64: input.CityId, Valid: input.CityId > 0},
		ImageUrl:    pgtype.Text{String: input.ImageUrl, Valid: input.ImageUrl != ""},
		HandoffMode: handoffMode,
		Address:     pgtype.Text{String: input.Address, Valid: input.Address != ""},
	}
	dto.Price.UnmarshalJSON([]byte(fmt.Sprintf("%f", input.Price)))
	dto.Weight.UnmarshalJSON([]byte(fmt.Sprintf("%f", input.Weight)))

	err = h.service.Update(pgtype.Int8{Int64: idInt, Valid: true}, dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logs.AddFromRequest(r, "Modification d'annonce", fmt.Sprintf("Annonce #%d modifiée: %s", idInt, dto.Name), "info")

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

	logs.AddFromRequest(r, "Suppression d'annonce", fmt.Sprintf("Annonce #%d: %s supprimée", idInt, existing.Name), "info")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "listing deleted successfully"}`)
}
