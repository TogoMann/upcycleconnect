package reporting

import (
	"backend/internal/middlewares"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /reporting/audit/items/pdf", middlewares.AdminOnly(handler.ExportAuditPDF))
	r.HandleFunc("GET /reporting/actors", middlewares.AdminOnly(handler.GetActorStats))
	r.HandleFunc("GET /reporting/prestations", middlewares.AdminOnly(handler.GetPrestationStats))
	r.HandleFunc("GET /reporting/predictions", middlewares.AdminOnly(handler.GetUserPredictions))
}
