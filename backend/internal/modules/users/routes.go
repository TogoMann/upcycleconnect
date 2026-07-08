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

	r.HandleFunc("GET /users", middlewares.AdminOnly(handler.GetAll))
	r.HandleFunc("GET /users/me", middlewares.Authenticated(handler.GetMe))
	r.HandleFunc("GET /users/{id}", middlewares.AdminOnly(handler.GetById))
	r.HandleFunc("GET /users/{id}/score", middlewares.Authenticated(handler.GetScore))
	r.HandleFunc("GET /users/{id}/score/history", middlewares.Authenticated(handler.GetScoreHistory))
	r.HandleFunc("GET /users/me/quests", middlewares.Authenticated(handler.GetMyQuests))
	r.HandleFunc("PATCH /users/{id}/tutorial", middlewares.Authenticated(handler.UpdateTutorialSeen))
	r.HandleFunc("POST /users", middlewares.AdminOnly(handler.Create))
	r.HandleFunc("PUT /users/{id}", middlewares.AdminOnly(handler.Update))
	r.HandleFunc("PATCH /users/{id}", middlewares.Authenticated(handler.Update))
	r.HandleFunc("DELETE /users/{id}", middlewares.AdminOnly(handler.DeleteById))
}
