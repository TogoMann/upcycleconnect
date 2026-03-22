package course

import (
	"net/http"
	
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /course", handler.GetAll)
	r.HandleFunc("GET /course/{id}", handler.GetById)
	r.HandleFunc("POST /course/", handler.Create)
	r.HandleFunc("DELETE /course/{id}", handler.DeleteById)
}
