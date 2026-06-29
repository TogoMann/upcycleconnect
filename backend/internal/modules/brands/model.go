package brands

import "github.com/jackc/pgx/v5/pgtype"

type Brand struct {
	Id          pgtype.Int8      `db:"id" json:"id"`
	Name        string           `db:"name" json:"name"`
	Description string           `db:"description" json:"description"`
	LogoUrl     string           `db:"logo_url" json:"logo_url"`
	Website     string           `db:"website" json:"website"`
	CreatedBy   pgtype.Int8      `db:"created_by" json:"created_by"`
	IsActive    bool             `db:"is_active" json:"is_active"`
	CreatedAt   pgtype.Timestamp `db:"created_at" json:"created_at"`
}
