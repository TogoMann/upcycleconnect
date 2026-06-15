package notifications

import (
	"backend/internal/middlewares"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /admin/notifications", middlewares.AdminOnly(handler.GetAll))
	r.HandleFunc("POST /admin/notifications", middlewares.AdminOnly(handler.Create))
}
