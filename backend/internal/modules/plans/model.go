package plans

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Plan struct {
	Id           pgtype.Int8      `db:"id" json:"id"`
	Name         string           `db:"name" json:"name"`
	Description  string           `db:"description" json:"description"`
	Price        pgtype.Numeric   `db:"price" json:"price"`
	BillingCycle string           `db:"billing_cycle" json:"billing_cycle"`
	Features     []string         `db:"features" json:"features"`
	IsActive     bool             `db:"is_active" json:"is_active"`
	CreatedAt    pgtype.Timestamp `db:"created_at" json:"created_at"`
}

type PlanFrontend struct {
	Id           int64    `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Price        float64  `json:"price"`
	BillingCycle string   `json:"billing_cycle"`
	Features     []string `json:"features"`
	IsActive     bool     `json:"is_active"`
}
