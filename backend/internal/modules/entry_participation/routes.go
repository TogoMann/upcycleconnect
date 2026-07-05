package entryparticipation

import (
	"backend/internal/middlewares"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /entry_participation", middlewares.AdminOnly(handler.GetAll))
	r.HandleFunc("GET /entry_participation/{id}", middlewares.Authenticated(handler.GetById))
	r.HandleFunc("GET /users/{user_id}/planning", middlewares.Authenticated(handler.GetByUser))
	r.HandleFunc("POST /entry_participation/", middlewares.Authenticated(handler.Create))
	r.HandleFunc("DELETE /entry_participation/{id}", middlewares.Authenticated(handler.DeleteById))
}
