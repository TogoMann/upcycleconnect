package users

import (
	"backend/internal/middlewares"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /users", handler.GetAll)
	r.HandleFunc("GET /users/me", middlewares.Authenticated(handler.GetMe))
	r.HandleFunc("GET /users/{id}", handler.GetById)
	r.HandleFunc("GET /users/{id}/score", handler.GetScore)
	r.HandleFunc("PATCH /users/{id}/tutorial", handler.UpdateTutorialSeen)

	r.HandleFunc("POST /users/", handler.Create)

	r.HandleFunc("PUT /users/{id}", middlewares.AdminOnly(handler.Update))
	r.HandleFunc("PATCH /users/{id}", middlewares.Authenticated(handler.Update))

	r.HandleFunc("DELETE /users/{id}", handler.DeleteById)
}
