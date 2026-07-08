package eventparticipation

import (
	"backend/internal/middlewares"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /event-participation/", middlewares.Authenticated(handler.GetAll))
	r.HandleFunc("GET /event-participation/{id}", middlewares.Authenticated(handler.GetById))
	r.HandleFunc("POST /event-participation/", middlewares.Authenticated(handler.Create))
	r.HandleFunc("DELETE /event-participation/{id}", middlewares.AdminOnly(handler.DeleteById))
}
