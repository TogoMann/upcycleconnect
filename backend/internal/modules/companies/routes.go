package companies

import (
	"backend/internal/middlewares"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	r.HandleFunc("GET /companies", middlewares.AdminOnly(handler.GetAll))
	r.HandleFunc("GET /companies/{id}", middlewares.AdminOnly(handler.GetById))
	r.HandleFunc("POST /companies", middlewares.AdminOnly(handler.Create))
	r.HandleFunc("GET /companies/siret", middlewares.AdminOnly(handler.GetBySiret))
}
