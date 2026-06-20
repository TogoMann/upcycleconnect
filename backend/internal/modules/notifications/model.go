package notifications

type Notification struct {
	Id      int64     `db:"id" json:"id"`
	Titre   string    `db:"titre" json:"titre"`
	Message string    `db:"message" json:"message"`
	Cible   string    `db:"cible" json:"cible"`
	Date    string    `db:"date" json:"date"`
	Envoyes int       `db:"envoyes" json:"envoyes"`
}

type NotificationRequest struct {
	Titre   string `json:"titre"`
	Message string `json:"message"`
	Cible   string `json:"cible"`
}
