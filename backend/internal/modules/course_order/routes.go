package courseorder

import (
	"net/http"
	
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /course_order", handler.GetAll)
	r.HandleFunc("GET /course_order/{id}", handler.GetById)
	r.HandleFunc("POST /course_order/", handler.Create)
	r.HandleFunc("DELETE /course_order/{id}", handler.DeleteById)
}
