package planning

import (
	"backend/internal/middlewares"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /planning/me", middlewares.Authenticated(handler.GetMyPlanning))
	r.HandleFunc("POST /planning/personal", middlewares.Authenticated(handler.CreatePersonalEvent))
	r.HandleFunc("DELETE /planning/personal/{id}", middlewares.Authenticated(handler.DeletePersonalEvent))
}
