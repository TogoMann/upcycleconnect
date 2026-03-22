package thread

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /thread", handler.GetAll)
	r.HandleFunc("GET /thread/{id}", handler.GetById)

	r.HandleFunc("POST /thread/", handler.Create)

	r.HandleFunc("DELETE /thread/{id}", handler.DeleteById)
}
