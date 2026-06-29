package brands

import (
	"backend/internal/middlewares"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /brands", handler.GetAll)
	r.HandleFunc("GET /brands/me", middlewares.Authenticated(handler.GetMyBrands))
	r.HandleFunc("GET /brands/{id}", handler.GetById)
	r.HandleFunc("POST /brands", middlewares.ProOnly(handler.Create))
	r.HandleFunc("PUT /brands/{id}", middlewares.ProOnly(handler.Update))
	r.HandleFunc("DELETE /brands/{id}", middlewares.ProOnly(handler.Delete))
	r.HandleFunc("GET /admin/brands", middlewares.AdminOnly(handler.GetAll))
	r.HandleFunc("DELETE /admin/brands/{id}", middlewares.AdminOnly(handler.Delete))
}
