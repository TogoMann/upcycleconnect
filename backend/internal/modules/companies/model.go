package companies

import "github.com/jackc/pgx/v5/pgtype"

type Company struct {
	Id        pgtype.Int8      `db:"id" json:"id"`
	Siret     string           `db:"siret" json:"siret"`
	Name      pgtype.Text      `db:"name" json:"name"`
	AddressId pgtype.Int8      `db:"address_id" json:"address_id"`
	CreatedAt pgtype.Timestamp `db:"created_at" json:"created_at"`
	Address   string           `json:"address,omitempty"`
}
