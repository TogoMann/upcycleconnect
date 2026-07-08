package settings

import (
	db "backend/internal/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Get() (*PlatformSettings, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT nom_site, logo_url, email_contact, telephone, adresse, commission_taux, maintenance, inscription_ouverte
		FROM platform_settings WHERE id = 1
	`)
	if err != nil {
		return nil, err
	}
	s, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[PlatformSettings])
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repository) Update(s PlatformSettings) error {
	_, err := r.db.Exec(db.Ctx, `
		UPDATE platform_settings SET
			nom_site = $1,
			logo_url = $2,
			email_contact = $3,
			telephone = $4,
			adresse = $5,
			commission_taux = $6,
			maintenance = $7,
			inscription_ouverte = $8
		WHERE id = 1
	`, s.NomSite, s.LogoUrl, s.EmailContact, s.Telephone, s.Adresse, s.CommissionTaux, s.Maintenance, s.InscriptionOuverte)
	return err
}

func (r *Repository) GetPublic() (*PublicSettings, error) {
	var s PublicSettings
	err := r.db.QueryRow(db.Ctx, `SELECT nom_site, maintenance FROM platform_settings WHERE id = 1`).Scan(&s.NomSite, &s.Maintenance)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repository) IsRegistrationOpen() (bool, error) {
	var open bool
	err := r.db.QueryRow(db.Ctx, `SELECT inscription_ouverte FROM platform_settings WHERE id = 1`).Scan(&open)
	if err != nil {
		return true, err
	}
	return open, nil
}
