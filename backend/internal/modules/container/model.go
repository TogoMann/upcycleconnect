package container

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type ConteneurFrontend struct {
	Id           int64  `db:"id" json:"id"`
	CodeBarres   string `db:"code_barres" json:"code_barres"`
	Localisation string `db:"localisation" json:"localisation"`
	Etat         string `db:"etat" json:"etat"`
	Capacite     int    `db:"capacite" json:"capacite"`
	Objets       int    `db:"objets" json:"objets"`
}

type Container struct {
	Id        pgtype.Int8      `db:"id" json:"id"`
	SiteId    pgtype.Int8      `db:"site_id" json:"site_id"`
	Status    string           `db:"status" json:"status"`
	CreatedAt pgtype.Timestamp `db:"created_at" json:"created_at"`
}

type Locker struct {
	Id          pgtype.Int8      `db:"id" json:"id"`
	ContainerId pgtype.Int8      `db:"container_id" json:"container_id"`
	Label       string           `db:"label" json:"label"`
	Status      string           `db:"status" json:"status"`
	Size        string           `db:"size" json:"size"`
	CreatedAt   pgtype.Timestamp `db:"created_at" json:"created_at"`
}
