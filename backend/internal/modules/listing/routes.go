package listing

import (
	"backend/internal/middlewares"
	"backend/internal/modules/users"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {

	repo := NewRepository(db)
	service := NewService(repo)

	userRepo := users.NewRepository(db)
	userService := users.NewService(userRepo)

	handler := NewHandler(service, userService)

	r.HandleFunc("GET /listing", handler.GetAllApproved)
	r.HandleFunc("GET /admin/listings", middlewares.AdminOnly(handler.GetAll))
	r.HandleFunc("GET /listing/{id}", handler.GetById)

	r.HandleFunc("POST /listing/", middlewares.Authenticated(handler.Create))

	r.HandleFunc("PUT /listing/{id}", middlewares.AdminOnly(handler.Update))
	r.HandleFunc("PATCH /listing/{id}", middlewares.AdminOnly(handler.Update))
	r.HandleFunc("PATCH /listing/{id}/approve", middlewares.AdminOnly(handler.Approve))
	r.HandleFunc("PATCH /listing/{id}/disapprove", middlewares.AdminOnly(handler.Disapprove))

	r.HandleFunc("DELETE /listing/{id}", handler.DeleteById)
}
