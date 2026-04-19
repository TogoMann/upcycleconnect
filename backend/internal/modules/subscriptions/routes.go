package subscriptions

import (
	"backend/internal/middlewares"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /admin/abonnements", middlewares.AdminOnly(handler.GetAllAbonnements))
	r.HandleFunc("GET /subscriptions", handler.GetAll)
	r.HandleFunc("GET /subscriptions/{id}", handler.GetById)
	r.HandleFunc("POST /subscriptions/", handler.Create)
	r.HandleFunc("PUT /subscriptions/{id}", middlewares.ProOnly(handler.Update))
	r.HandleFunc("PATCH /subscriptions/{id}", middlewares.ProOnly(handler.Update))
	r.HandleFunc("DELETE /subscriptions/{id}", middlewares.AdminOnly(handler.DeleteById))
}
