package entryparticipation

import "time"

type EntryStatus string

const (
	Accepted EntryStatus = "accepted"
	Declined EntryStatus = "declined"
	Pending  EntryStatus = "pending"
)

type EntryParticipation struct {
	Id       int64       `db:"id" json:"id"`
	EntryId  int64       `db:"entry_id" json:"entry_id"`
	UserId   int64       `db:"user_id" json:"user_id"`
	Status   EntryStatus `db:"status" json:"status"`
	JoinedAt time.Time   `db:"joined_at" json:"joined_at"`
}
