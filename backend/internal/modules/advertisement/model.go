package advertisement

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type PubFrontend struct {
	Id        int64   `db:"id" json:"id"`
	Titre     string  `db:"titre" json:"titre"`
	Annonceur string  `db:"annonceur" json:"annonceur"`
	Type      string  `db:"type" json:"type"`
	Debut     string  `db:"debut" json:"debut"`
	Fin       string  `db:"fin" json:"fin"`
	Statut    string  `db:"statut" json:"statut"`
	Budget    float64 `db:"budget" json:"budget"`
}

type Advertisement struct {
	Id                    pgtype.Int8      `db:"id" json:"id"`
	AnnouncerId           pgtype.Int8      `db:"announcer_id" json:"announcer_id"`
	TargetId              int64            `db:"target_id" json:"target_id"`
	TargetType            string           `db:"target_type" json:"target_type"`
	AdType                string           `db:"ad_type" json:"ad_type"`
	Budget                pgtype.Numeric   `db:"budget" json:"budget"`
	Status                string           `db:"status" json:"status"`
	StripePaymentIntentId string           `db:"stripe_payment_intent_id" json:"stripe_payment_intent_id"`
	StartDate             pgtype.Date      `db:"start_date" json:"start_date"`
	EndDate               pgtype.Date      `db:"end_date" json:"end_date"`
	CreatedAt             pgtype.Timestamp `db:"created_at" json:"created_at"`
	ApprovedBy            pgtype.Int8      `db:"approved_by" json:"approved_by"`
}
