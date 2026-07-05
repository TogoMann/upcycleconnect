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
			TO_CHAR(CURRENT_DATE, 'TMMonth YYYY') as periode
		FROM listing_order
		WHERE (status = 'paid' OR status = 'completed')
		AND date_trunc('month', created_at) = date_trunc('month', CURRENT_DATE)
	`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Commission])
}

func (r *Repository) GetFinancier() (*FinancierData, error) {
	var data FinancierData

	err := r.db.QueryRow(db.Ctx, "SELECT CAST(COALESCE(SUM(price), 0) AS FLOAT8) FROM listing_order WHERE status = 'paid' OR status = 'completed'").Scan(&data.CaTotal)
	if err != nil {
		return nil, err
	}

	var courseCa float64
	err = r.db.QueryRow(db.Ctx, "SELECT CAST(COALESCE(SUM(price), 0) AS FLOAT8) FROM course_order").Scan(&courseCa)
	if err != nil {
		return nil, err
	}

	data.CaTotal += courseCa

	err = r.db.QueryRow(db.Ctx, `
		SELECT CAST(COALESCE(SUM(price), 0) AS FLOAT8) FROM (
			SELECT price, created_at FROM listing_order WHERE status = 'paid' OR status = 'completed'
			UNION ALL
			SELECT price, booked_at as created_at FROM course_order
		) t
		WHERE date_trunc('month', created_at) = date_trunc('month', CURRENT_DATE)
	`).Scan(&data.CaMois)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(db.Ctx, "SELECT CAST(COALESCE(SUM(amount), 0) AS FLOAT8) FROM expenses").Scan(&data.Charges)
	if err != nil {
		return nil, err
	}

	data.Marge = data.CaTotal - data.Charges

	rows, err := r.db.Query(db.Ctx, `
		SELECT mois, CAST(COALESCE(SUM(ca), 0) AS FLOAT8) as ca
		FROM (
			SELECT TO_CHAR(date_trunc('month', created_at), 'TMMon') as mois, date_trunc('month', created_at) as m, price as ca
			FROM listing_order WHERE status = 'paid' OR status = 'completed'
			UNION ALL
			SELECT TO_CHAR(date_trunc('month', booked_at), 'TMMon') as mois, date_trunc('month', booked_at) as m, price as ca
			FROM course_order
		) t
		WHERE m >= date_trunc('month', CURRENT_DATE) - INTERVAL '3 months'
		GROUP BY mois, m
		ORDER BY m
	`)
	if err != nil {
		return nil, err
	}

	revenueByMonth := map[string]float64{}
	var monthOrder []string
	for rows.Next() {
		var mois string
		var ca float64
		if err := rows.Scan(&mois, &ca); err != nil {
			return nil, err
		}
		revenueByMonth[mois] = ca
		monthOrder = append(monthOrder, mois)
	}

	expenseRows, err := r.db.Query(db.Ctx, `
		SELECT TO_CHAR(date_trunc('month', created_at), 'TMMon') as mois, CAST(COALESCE(SUM(amount), 0) AS FLOAT8) as charges
		FROM expenses
		WHERE date_trunc('month', created_at) >= date_trunc('month', CURRENT_DATE) - INTERVAL '3 months'
		GROUP BY mois
	`)
	if err != nil {
		return nil, err
	}

	expensesByMonth := map[string]float64{}
	for expenseRows.Next() {
		var mois string
		var charges float64
		if err := expenseRows.Scan(&mois, &charges); err != nil {
			return nil, err
		}
		expensesByMonth[mois] = charges
	}

	data.Evolution = make([]Evolution, 0, len(monthOrder))
	for _, mois := range monthOrder {
		data.Evolution = append(data.Evolution, Evolution{
			Mois:    mois,
			Ca:      revenueByMonth[mois],
			Charges: expensesByMonth[mois],
		})
	}

	return &data, nil
}

func (r *Repository) GetReport() (*FinancialReport, error) {
	var report FinancialReport

	err := r.db.QueryRow(db.Ctx, "SELECT CAST(COALESCE(SUM(price), 0) AS FLOAT8) FROM listing_order WHERE status = 'paid' OR status = 'completed'").Scan(&report.ListingRevenue)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(db.Ctx, "SELECT CAST(COALESCE(SUM(price), 0) AS FLOAT8) FROM course_order").Scan(&report.CourseRevenue)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(db.Ctx, "SELECT CAST(COALESCE(SUM(price), 0) AS FLOAT8) FROM subscriptions").Scan(&report.SubscriptionRevenue)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(db.Ctx, "SELECT CAST(COALESCE(SUM(budget), 0) AS FLOAT8) FROM advertisement WHERE status = 'validated'").Scan(&report.AdRevenue)
	if err != nil {
		return nil, err
	}

	report.TotalCommissions = report.ListingRevenue * 0.10

	report.TotalRevenue = report.ListingRevenue + report.CourseRevenue + report.SubscriptionRevenue + report.AdRevenue

	return &report, nil
}

func (r *Repository) CreateExpense(e Expense) error {
	_, err := r.db.Exec(db.Ctx,
		"INSERT INTO expenses (label, amount, category, created_by) VALUES ($1, $2, $3, $4)",
		e.Label, e.Amount, e.Category, e.CreatedBy)
	return err
}

func (r *Repository) GetAllExpenses() ([]Expense, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT id, label, CAST(amount AS FLOAT8) as amount, COALESCE(category, '') as category, created_by,
		TO_CHAR(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as created_at
		FROM expenses ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Expense])
}

func (r *Repository) CreateInvoice(inv Invoice) error {
	_, err := r.db.Exec(db.Ctx, `
		INSERT INTO invoices (user_id, seller_id, order_id, order_type, invoice_number, amount, vat_amount, total_amount, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, inv.UserId, inv.SellerId, inv.OrderId, inv.OrderType, inv.InvoiceNumber, inv.Amount, inv.VatAmount, inv.TotalAmount, inv.Status)
	return err
}

func (r *Repository) GetUserInvoices(userId int64) ([]Invoice, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT id, user_id, seller_id, order_id, order_type, invoice_number, amount, vat_amount, total_amount, status, 
		TO_CHAR(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as created_at
		FROM invoices
		WHERE user_id = $1 OR seller_id = $1
		ORDER BY created_at DESC
	`, userId)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Invoice])
}

func (r *Repository) GetInvoiceById(id int64) (*InvoiceDetail, error) {
	var detail InvoiceDetail
	err := r.db.QueryRow(db.Ctx, `
		SELECT i.id, i.user_id, i.seller_id, i.order_id, i.order_type, i.invoice_number, i.amount, i.vat_amount, i.total_amount, i.status,
		TO_CHAR(i.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as created_at,
		u.first_name || ' ' || u.last_name as user_name,
		u.email as user_email,
		COALESCE(l.name, c.name, CASE WHEN i.order_type = 'event' THEN 'Participation événement' END, s.tier, 'Transaction') as description
		FROM invoices i
		JOIN users u ON i.user_id = u.id
		LEFT JOIN listing_order lo ON i.order_id = lo.id AND i.order_type = 'listing'
		LEFT JOIN listing l ON lo.listing_id = l.id
		LEFT JOIN course_order co ON i.order_id = co.id AND i.order_type = 'course'
		LEFT JOIN course c ON co.course_id = c.id
		LEFT JOIN subscriptions s ON i.order_id = s.id AND i.order_type = 'subscription'
		WHERE i.id = $1
	`, id).Scan(
		&detail.Id, &detail.UserId, &detail.SellerId, &detail.OrderId, &detail.OrderType, &detail.InvoiceNumber, &detail.Amount, &detail.VatAmount, &detail.TotalAmount, &detail.Status, &detail.CreatedAt,
		&detail.UserName, &detail.UserEmail, &detail.Description,
	)
	if err != nil {
		return nil, err
	}
	return &detail, nil
}
