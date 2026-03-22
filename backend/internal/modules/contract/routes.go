package contract

import (
	"net/http"
	
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /contract", handler.GetAll)
	r.HandleFunc("GET /contract/{id}", handler.GetById)
	r.HandleFunc("POST /contract/", handler.Create)
	r.HandleFunc("DELETE /contract/{id}", handler.DeleteById)
}
