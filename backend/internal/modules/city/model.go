package city

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type City struct {
	Id      pgtype.Int8 `db:"id" json:"id"`
	Name    string      `db:"name" json:"name"`
	ZipCode string      `db:"zip_code" json:"zip_code"`
}
