package item

import (
	"backend/internal/middlewares"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /items/me", middlewares.Authenticated(handler.GetByUserId))
	r.HandleFunc("GET /items", handler.GetAll)
	r.HandleFunc("GET /items/{id}", handler.GetById)
	r.HandleFunc("POST /items/", middlewares.Authenticated(handler.Create))
	r.HandleFunc("PUT /items/{id}", middlewares.AdminOnly(handler.Update))
	r.HandleFunc("DELETE /items/{id}", middlewares.AdminOnly(handler.Delete))
	r.HandleFunc("POST /items/{id}/collect", middlewares.ProOnly(handler.Collect))
}
