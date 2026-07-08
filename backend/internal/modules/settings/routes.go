package settings

import (
	"backend/internal/middlewares"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) *Service {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /admin/parametres", middlewares.AdminOnly(handler.Get))
	r.HandleFunc("PUT /admin/parametres", middlewares.AdminOnly(handler.Update))
	r.HandleFunc("GET /parametres/public", handler.GetPublic)

	return service
}
