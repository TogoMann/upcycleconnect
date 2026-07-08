package financial

import (
	"backend/internal/middlewares"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /admin/financier", middlewares.AdminOnly(handler.GetFinancier))
	r.HandleFunc("GET /admin/commissions", middlewares.AdminOnly(handler.GetCommissions))
	r.HandleFunc("GET /admin/financial/report", middlewares.AdminOnly(handler.GetReport))
	r.HandleFunc("GET /admin/expenses", middlewares.AdminOnly(handler.GetAllExpenses))
	r.HandleFunc("POST /admin/expenses", middlewares.AdminOnly(handler.CreateExpense))
	r.HandleFunc("GET /financial/my-invoices", middlewares.Authenticated(handler.GetMyInvoices))
	r.HandleFunc("GET /financial/invoices/{id}/pdf", middlewares.Authenticated(handler.ExportInvoicePDF))
}
