package listingorder

import (
	"backend/internal/middlewares"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /listing-order/me", middlewares.Authenticated(handler.GetByUserId))
	r.HandleFunc("GET /listing-order", middlewares.AdminOnly(handler.GetAll))
	r.HandleFunc("GET /listing-order/{id}", middlewares.Authenticated(handler.GetById))

	r.HandleFunc("POST /listing-order/", middlewares.Authenticated(handler.Create))

	r.HandleFunc("DELETE /listing-order/{id}", middlewares.AdminOnly(handler.DeleteById))
}
