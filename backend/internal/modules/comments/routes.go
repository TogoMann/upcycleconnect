package comments

import (
	"net/http"
	
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /comments", handler.GetAll)
	r.HandleFunc("GET /comments/{id}", handler.GetById)
	r.HandleFunc("POST /comments/", handler.Create)
	r.HandleFunc("DELETE /comments/{id}", handler.DeleteById)
}
