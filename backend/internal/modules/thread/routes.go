package thread

import (
	"net/http"

	"github.com/jackc/pgx/v5"
)

func RegisterRoutes(r *http.ServeMux, db *pgx.Conn) {

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /routes", handler.GetAll)
	r.HandleFunc("GET /routes/{id}", handler.GetById)

	r.HandleFunc("POST /routes/", handler.Create)

	r.HandleFunc("DELETE /routes/{id}", handler.DeleteById)
}
