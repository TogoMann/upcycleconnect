package subscriptions

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type AbonnementFrontend struct {
	Id          int64   `db:"id" json:"id"`
	Entreprise  string  `db:"entreprise" json:"entreprise"`
	Utilisateur string  `db:"utilisateur" json:"utilisateur"`
	Plan        string  `db:"plan" json:"plan"`
	Debut       string  `db:"debut" json:"debut"`
	Fin         string  `db:"fin" json:"fin"`
	Statut      string  `db:"statut" json:"statut"`
	Montant     float64 `db:"montant" json:"montant"`
}

type Subscription struct {
	Id           pgtype.Int8 `db:"id" json:"id"`
	SubscriberId pgtype.Int8 `db:"subscriber_id" json:"subscriber_id"`
	Price        float64     `db:"price" json:"price"`
	Tier         string      `db:"tier" json:"tier"`
	CreatedAt    pgtype.Date `db:"created_at" json:"created_at"`
	Until        pgtype.Date `db:"until" json:"until"`
}
