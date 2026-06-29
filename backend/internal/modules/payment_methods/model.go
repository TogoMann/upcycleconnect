package payment_methods

import "github.com/jackc/pgx/v5/pgtype"

type PaymentMethod struct {
	Id                    pgtype.Int8      `db:"id" json:"id"`
	UserId                pgtype.Int8      `db:"user_id" json:"user_id"`
	StripePaymentMethodId string           `db:"stripe_payment_method_id" json:"stripe_payment_method_id"`
	CardLast4             string           `db:"card_last4" json:"card_last4"`
	CardBrand             string           `db:"card_brand" json:"card_brand"`
	CardExpMonth          int32            `db:"card_exp_month" json:"card_exp_month"`
	CardExpYear           int32            `db:"card_exp_year" json:"card_exp_year"`
	IsDefault             bool             `db:"is_default" json:"is_default"`
	CreatedAt             pgtype.Timestamp `db:"created_at" json:"created_at"`
}
