package listingorder

import (
	"net/http"

	
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /listing-order", handler.GetAll)
	r.HandleFunc("GET /listing-order/{id}", handler.GetById)

	r.HandleFunc("POST /listing-order/", handler.Create)

	r.HandleFunc("DELETE /listing-order/{id}", handler.DeleteById)
}
