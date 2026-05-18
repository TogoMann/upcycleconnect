package companies

import "github.com/jackc/pgx/v5/pgtype"

type Company struct {
	Id        pgtype.Int8      `db:"id" json:"id"`
	Siret     string           `db:"siret" json:"siret"`
	Name      pgtype.Text      `db:"name" json:"name"`
	Address   pgtype.Text      `db:"address" json:"address"`
	CreatedAt pgtype.Timestamp `db:"created_at" json:"created_at"`
}
