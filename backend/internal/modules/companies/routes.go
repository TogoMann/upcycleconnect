package companies

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	r.HandleFunc("GET /companies", handler.GetAll)
	r.HandleFunc("GET /companies/{id}", handler.GetById)
	r.HandleFunc("POST /companies", handler.Create)
	r.HandleFunc("GET /companies/siret", handler.GetBySiret)
}
