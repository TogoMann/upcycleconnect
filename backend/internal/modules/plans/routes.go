package plans

import (
	"backend/internal/middlewares"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /plans", handler.GetAll)
	r.HandleFunc("GET /plans/{id}", handler.GetById)
	r.HandleFunc("POST /admin/plans", middlewares.AdminOnly(handler.Create))
	r.HandleFunc("PUT /admin/plans/{id}", middlewares.AdminOnly(handler.Update))
	r.HandleFunc("DELETE /admin/plans/{id}", middlewares.AdminOnly(handler.Delete))
}
