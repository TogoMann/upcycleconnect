package users

import (
	"net/http"

	"github.com/jackc/pgx/v5"
)

func RegisterRoutes(r *http.ServeMux, db *pgx.Conn) {

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /users", handler.GetAll)
	r.HandleFunc("GET /users/{id}", handler.GetById)

	r.HandleFunc("POST /users/", handler.Create)

	r.HandleFunc("DELETE /users/{id}", handler.DeleteById)
}
