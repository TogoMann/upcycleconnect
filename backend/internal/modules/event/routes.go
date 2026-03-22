package event

import (
	"net/http"
	
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /event", handler.GetAll)
	r.HandleFunc("GET /event/{id}", handler.GetById)
	r.HandleFunc("POST /event/", handler.Create)
	r.HandleFunc("DELETE /event/{id}", handler.DeleteById)
}
