package eventparticipation

type Event struct {
	EventId int64 `db:"event_id" json:"event_id"`
	UserId  int64 `db:"user_id" json:"user_id"`
}
