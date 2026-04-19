package entry

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type DepotFrontend struct {
	Id          int64  `db:"id" json:"id"`
	Utilisateur string `db:"utilisateur" json:"utilisateur"`
	Objet       string `db:"objet" json:"objet"`
	Date        string `db:"date" json:"date"`
	Statut      string `db:"statut" json:"statut"`
	CodeEnvoye  bool   `db:"code_envoye" json:"code_envoye"`
}

type Entry struct {
	Id        pgtype.Int8      `db:"id" json:"id"`
	CreatedBy pgtype.Int8      `db:"created_by" json:"created_by"`
	CreatedAt pgtype.Timestamp `db:"created_at" json:"created_at"`
	Schedule  pgtype.Date      `db:"schedule" json:"schedule"`
	Start     pgtype.Time      `db:"start" json:"start"`
	Ending    pgtype.Time      `db:"ending" json:"ending"`
}
