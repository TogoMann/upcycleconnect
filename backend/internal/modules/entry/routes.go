package entry

import (
	"backend/internal/middlewares"
	"net/http"
	
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /admin/depots", middlewares.AdminOnly(handler.GetAllDepots))
	r.HandleFunc("POST /admin/depots/{id}/valider", middlewares.AdminOnly(handler.ValiderDepot))
	r.HandleFunc("POST /admin/depots/{id}/code", middlewares.AdminOnly(handler.EnvoyerCode))
	r.HandleFunc("GET /entry", middlewares.AdminOnly(handler.GetAll))
	r.HandleFunc("GET /entry/{id}", middlewares.AdminOnly(handler.GetById))
	r.HandleFunc("POST /entry", handler.Create)
	r.HandleFunc("DELETE /entry/{id}", middlewares.AdminOnly(handler.DeleteById))
}
