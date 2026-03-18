package listing

import (
	"net/http"

	"github.com/jackc/pgx/v5"
)

func RegisterRoutes(r *http.ServeMux, db *pgx.Conn) {

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /listing", handler.GetAll)
	r.HandleFunc("GET /listing/{id}", handler.GetById)

	r.HandleFunc("POST /listing/", handler.Create)

	r.HandleFunc("DELETE /listing/{id}", handler.DeleteById)
}
