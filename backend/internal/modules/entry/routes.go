package entry

import (
	"net/http"
	"github.com/jackc/pgx/v5"
)

func RegisterRoutes(r *http.ServeMux, db *pgx.Conn) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /entry", handler.GetAll)
	r.HandleFunc("GET /entry/{id}", handler.GetById)
	r.HandleFunc("POST /entry/", handler.Create)
	r.HandleFunc("DELETE /entry/{id}", handler.DeleteById)
}
