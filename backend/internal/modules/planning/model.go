package planning

import "github.com/jackc/pgx/v5/pgtype"

type PlanningItem struct {
	Id          int64  `json:"id"`
	Type        string `json:"type"` // "depot", "workshop", "event", "personal"
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Location    string `json:"location"`
}

type AdminPlanningItem struct {
	Id           int64  `db:"id" json:"id"`
	Titre        string `db:"titre" json:"titre"`
	Type         string `db:"type" json:"type"`
	Responsable  string `db:"responsable" json:"responsable"`
	Date         string `db:"date" json:"date"`
	HeureDebut   string `db:"heure_debut" json:"heure_debut"`
	HeureFin     string `db:"heure_fin" json:"heure_fin"`
	Participants int    `db:"participants" json:"participants"`
}

type PersonalEvent struct {
	Id          pgtype.Int8      `db:"id" json:"id"`
	UserId      pgtype.Int8      `db:"user_id" json:"user_id"`
	Title       string           `db:"title" json:"title"`
	Description string           `db:"description" json:"description"`
	Date        pgtype.Date      `db:"date" json:"date"`
	StartTime   pgtype.Time      `db:"start_time" json:"start_time"`
	EndTime     pgtype.Time      `db:"end_time" json:"end_time"`
	CreatedAt   pgtype.Timestamp `db:"created_at" json:"created_at"`
}
