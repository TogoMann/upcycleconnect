package entry

import "time"

type Entry struct {
	Id        int64     `db:"id" json:"id"`
	CreatedBy int64     `db:"created_by" json:"created_by"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Schedule  time.Time `db:"schedule" json:"schedule"`
	Start     time.Time `db:"start" json:"start"`
	Ending    time.Time `db:"ending" json:"ending"`
}
