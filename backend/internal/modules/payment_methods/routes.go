package payment_methods

import (
	"backend/internal/middlewares"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /payment-methods/me", middlewares.Authenticated(handler.GetMyPaymentMethods))
	r.HandleFunc("GET /payment-methods/check", middlewares.Authenticated(handler.HasPaymentMethod))
	r.HandleFunc("POST /payment-methods", middlewares.Authenticated(handler.Create))
	r.HandleFunc("DELETE /payment-methods/{id}", middlewares.Authenticated(handler.Delete))
	r.HandleFunc("PATCH /payment-methods/{id}/default", middlewares.Authenticated(handler.SetDefault))
	r.HandleFunc("GET /admin/payment-methods", middlewares.AdminOnly(handler.GetAll))
}
