package settings

type PlatformSettings struct {
	NomSite            string  `db:"nom_site" json:"nom_site"`
	LogoUrl            string  `db:"logo_url" json:"logo_url"`
	EmailContact       string  `db:"email_contact" json:"email_contact"`
	Telephone          string  `db:"telephone" json:"telephone"`
	Adresse            string  `db:"adresse" json:"adresse"`
	CommissionTaux     float64 `db:"commission_taux" json:"commission_taux"`
	Maintenance        bool    `db:"maintenance" json:"maintenance"`
	InscriptionOuverte bool    `db:"inscription_ouverte" json:"inscription_ouverte"`
}

type PublicSettings struct {
	NomSite     string `json:"nom_site"`
	Maintenance bool   `json:"maintenance"`
}
