package eventparticipation

import "github.com/jackc/pgx/v5/pgtype"

type Event struct {
	EventId pgtype.Int8 `db:"event_id" json:"event_id"`
	UserId  pgtype.Int8 `db:"user_id" json:"user_id"`
}
