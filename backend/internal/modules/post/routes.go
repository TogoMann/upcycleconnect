package post

import (
	"net/http"

	"github.com/jackc/pgx/v5"
)

func RegisterRoutes(r *http.ServeMux, db *pgx.Conn) {

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /post", handler.GetAll)
	r.HandleFunc("GET /post/{id}", handler.GetById)

	r.HandleFunc("POST /post/", handler.Create)
	// TODO: UPDATE post

	r.HandleFunc("DELETE /post/{id}", handler.DeleteById)
}
