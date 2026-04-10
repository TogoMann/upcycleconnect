package project

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /project", handler.GetAll)
	r.HandleFunc("GET /project/{id}", handler.GetById)
	r.HandleFunc("GET /project/{id}/steps", handler.GetSteps)
	r.HandleFunc("POST /project/", handler.Create)
	r.HandleFunc("DELETE /project/{id}", handler.DeleteById)
}
