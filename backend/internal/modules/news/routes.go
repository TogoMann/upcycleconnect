package news

import (
	"net/http"

	"github.com/jackc/pgx/v5"
)

func RegisterRoutes(r *http.ServeMux, db *pgx.Conn) {

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /news", handler.GetAll)
	r.HandleFunc("GET /news/{id}", handler.GetById)

	r.HandleFunc("POST /news/", handler.Create)

	r.HandleFunc("DELETE /news/{id}", handler.DeleteById)
}
