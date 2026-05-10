package listing

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type ListingStatus string

const (
	Active    ListingStatus = "active"
	Sold      ListingStatus = "sold"
	Cancelled ListingStatus = "cancelled"
)

type ListingCategory string

const (
	Mobilier     ListingCategory = "Mobilier"
	Decoration   ListingCategory = "Décoration"
	Vetements    ListingCategory = "Vêtements"
	Jouet        ListingCategory = "Jouet"
	Electronique ListingCategory = "Electronique"
	Outils       ListingCategory = "Outils"
)

type Listing struct {
	Id          pgtype.Int8      `db:"id" json:"id"`
	Name        string           `db:"name" json:"name"`
	Description string           `db:"description" json:"description"`
	Category    ListingCategory  `db:"category" json:"category"`
	ItemId      pgtype.Int8      `db:"item_id" json:"item_id"`
	CityId      pgtype.Int8      `db:"city_id" json:"city_id"`
	CityName      string           `db:"city_name" json:"city_name"`
	CreatedByName string           `db:"created_by_name" json:"created_by_name"`
	CreatedBy     pgtype.Int8      `db:"created_by" json:"created_by"`
	CreatedAt   pgtype.Timestamp `db:"created_at" json:"created_at"`
	Approved    bool             `db:"approved" json:"approved"`
	ApprovedBy  pgtype.Int8      `db:"approved_by" json:"approved_by"`
	ApprovedAt  pgtype.Timestamp `db:"approved_at" json:"approved_at"`
	Status      ListingStatus    `db:"status" json:"status"`
	Price       pgtype.Numeric   `db:"price" json:"price"`
	ImageUrl    pgtype.Text      `db:"image_url" json:"image_url"`
}
