package logs

type Log struct {
	Id          int64  `db:"id" json:"id"`
	Utilisateur string `db:"utilisateur" json:"utilisateur"`
	Action      string `db:"action" json:"action"`
	Ressource   string `db:"ressource" json:"ressource"`
	Ip          string `db:"ip" json:"ip"`
	Date        string `json:"date"`
	Niveau      string `db:"niveau" json:"niveau"`
}
