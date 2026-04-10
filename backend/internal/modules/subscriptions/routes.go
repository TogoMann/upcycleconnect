package subscriptions

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /subscriptions", handler.GetAll)
	r.HandleFunc("GET /subscriptions/{id}", handler.GetById)
	r.HandleFunc("POST /subscriptions/", handler.Create)
	r.HandleFunc("DELETE /subscriptions/{id}", handler.DeleteById)
}
