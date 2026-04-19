package advertisement

import (
	"backend/internal/middlewares"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /admin/publicites", middlewares.AdminOnly(handler.GetAllPubs))
	r.HandleFunc("PATCH /admin/publicites/{id}", middlewares.AdminOnly(handler.UpdateStatus))
	r.HandleFunc("DELETE /admin/publicites/{id}", middlewares.AdminOnly(handler.Delete))

	r.HandleFunc("GET /advertisement", handler.GetAll)
	r.HandleFunc("GET /advertisement/{id}", handler.GetById)
	r.HandleFunc("POST /advertisement", handler.Create)
	r.HandleFunc("PATCH /advertisement/{id}/approve", middlewares.AdminOnly(handler.Approve))
	r.HandleFunc("PATCH /advertisement/{id}/reject", middlewares.AdminOnly(handler.Reject))
	r.HandleFunc("DELETE /advertisement/{id}", middlewares.AdminOnly(handler.Delete))
}
