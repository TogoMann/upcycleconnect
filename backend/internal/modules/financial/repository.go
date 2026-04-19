package financial

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

func (r *Repository) GetCommissions() ([]Commission, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			1 as id,
			'Ventes catalogue' as type,
			10 as taux,
			CAST(COALESCE(SUM(price), 0) * 0.10 AS FLOAT8) as montant_total,
			CAST(COUNT(*) AS INTEGER) as nb_transactions,
			'Avril 2026' as periode
		FROM listing_order
		WHERE status = 'paid' OR status = 'completed'
	`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Commission])
}

func (r *Repository) GetFinancier() (*FinancierData, error) {
	var data FinancierData

	err := r.db.QueryRow(db.Ctx, "SELECT CAST(COALESCE(SUM(price), 0) AS FLOAT8) FROM listing_order WHERE status = 'paid' OR status = 'completed'").Scan(&data.CaTotal)
	if err != nil { return nil, err }
	
	var courseCa float64
	err = r.db.QueryRow(db.Ctx, "SELECT CAST(COALESCE(SUM(price), 0) AS FLOAT8) FROM course_order").Scan(&courseCa)
	if err != nil { return nil, err }
	
	data.CaTotal += courseCa

	data.CaMois = data.CaTotal * 0.15 
	data.Charges = data.CaTotal * 0.4
	data.Marge = data.CaTotal - data.Charges

	data.Evolution = []Evolution{
		{Mois: "Jan", Ca: 1200, Charges: 500},
		{Mois: "Feb", Ca: 1500, Charges: 600},
		{Mois: "Mar", Ca: 1100, Charges: 450},
		{Mois: "Apr", Ca: data.CaTotal, Charges: data.Charges},
	}

	return &data, nil
}

func (r *Repository) GetReport() (*FinancialReport, error) {
	var report FinancialReport

	err := r.db.QueryRow(db.Ctx, "SELECT CAST(COALESCE(SUM(price), 0) AS FLOAT8) FROM listing_order WHERE status = 'paid' OR status = 'completed'").Scan(&report.ListingRevenue)
	if err != nil { return nil, err }
	
	err = r.db.QueryRow(db.Ctx, "SELECT CAST(COALESCE(SUM(price), 0) AS FLOAT8) FROM course_order").Scan(&report.CourseRevenue)
	if err != nil { return nil, err }
	
	err = r.db.QueryRow(db.Ctx, "SELECT CAST(COALESCE(SUM(price), 0) AS FLOAT8) FROM subscriptions").Scan(&report.SubscriptionRevenue)
	if err != nil { return nil, err }
	
	err = r.db.QueryRow(db.Ctx, "SELECT CAST(COALESCE(SUM(budget), 0) AS FLOAT8) FROM advertisement WHERE status = 'validated'").Scan(&report.AdRevenue)
	if err != nil { return nil, err }

	report.TotalCommissions = report.ListingRevenue * 0.10
	
	report.TotalRevenue = report.ListingRevenue + report.CourseRevenue + report.SubscriptionRevenue + report.AdRevenue

	return &report, nil
}
