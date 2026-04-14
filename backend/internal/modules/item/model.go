package item

import "github.com/jackc/pgx/v5/pgtype"

type ItemStatus string

const (
	Deposited ItemStatus = "deposited"
	Validated ItemStatus = "validated"
	Collected ItemStatus = "collected"
)

type ItemState string

const (
	Neuf    ItemState = "Neuf"
	BonEtat ItemState = "Bon état"
	Abime   ItemState = "Abimé"
	Casse   ItemState = "Cassé"
)

type Item struct {
	Id            pgtype.Int8      `db:"id" json:"id"`
	OwnerId       pgtype.Int8      `db:"owner_id" json:"owner_id"`
	ContainerId   pgtype.Int8      `db:"container_id" json:"container_id"`
	SiteId        pgtype.Int8      `db:"site_id" json:"site_id"`
	MaterialType  string           `db:"material_type" json:"material_type"`
	PhysicalState ItemState        `db:"physical_state" json:"physical_state"`
	Status        ItemStatus       `db:"status" json:"status"`
	Weight        pgtype.Numeric   `db:"weight" json:"weight"`
	CreatedAt     pgtype.Timestamp `db:"created_at" json:"created_at"`
}
