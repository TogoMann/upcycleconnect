package subscriptions

import "github.com/jackc/pgx/v5/pgtype"

type Subscriptions struct {
	Id           pgtype.Int8    `db:"id" json:"id"`
	SubscriberId pgtype.Int8    `db:"subscriber_id" json:"subscriber_id"`
	Price        pgtype.Numeric `db:"price" json:"price"`
	Tier         string         `db:"tier" json:"tier"`
	CreatedAt    pgtype.Date    `db:"created_at" json:"created_at"`
	Until        pgtype.Date    `db:"until" json:"until"`
}
