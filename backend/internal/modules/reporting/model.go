package reporting

import "time"

type AuditReport struct {
	Title       string    `json:"title"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	ItemCount   int       `json:"item_count"`
	GeneratedAt time.Time `json:"generated_at"`
}

type ActorStats struct {
	Role  string `json:"role"`
	Count int    `json:"count"`
}

type PrestationStats struct {
	Type  string  `json:"type"`
	Count int     `json:"count"`
	Total float64 `json:"total_revenue"`
}

type UserPrediction struct {
	UserID               int64     `json:"user_id"`
	Username             string    `json:"username"`
	Email                string    `json:"email"`
	PredictedServiceType string    `json:"predicted_service_type"`
	Probability          float64   `json:"probability"`
	CalculatedAt         time.Time `json:"calculated_at"`
}

type PaginatedPredictions struct {
	Data       []UserPrediction `json:"data"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
}
