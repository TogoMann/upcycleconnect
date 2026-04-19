package course

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type OffreFrontend struct {
	Id          int64   `db:"id" json:"id"`
	Nom         string  `db:"nom" json:"nom"`
	Categorie   string  `db:"categorie" json:"categorie"`
	Prix        float64 `db:"prix" json:"prix"`
	Description string  `db:"description" json:"description"`
	Actif       bool    `db:"actif" json:"actif"`
}

type Course struct {
	Id          pgtype.Int8      `db:"id" json:"id"`
	Name        string           `db:"name" json:"name"`
	Description string           `db:"description" json:"description"`
	MaxCapacity pgtype.Int4      `db:"max_capacity" json:"max_capacity"`
	CreatedBy   pgtype.Int8      `db:"created_by" json:"created_by"`
	CreatedAt   pgtype.Timestamp `db:"created_at" json:"created_at"`
	Approved    bool             `db:"approved" json:"approved"`
	ApprovedBy  pgtype.Int8      `db:"approved_by" json:"approved_by"`
	ApprovedAt  pgtype.Timestamp `db:"approved_at" json:"approved_at"`
	Price       pgtype.Numeric   `db:"price" json:"price"`
}

type UserCourse struct {
	Course
	BuyerID  pgtype.Int8      `db:"buyer_id" json:"buyer_id"`
	BookedAt pgtype.Timestamp `db:"booked_at" json:"booked_at"`
}
