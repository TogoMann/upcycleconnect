package event

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

	r.HandleFunc("GET /event", handler.GetAll)
	r.HandleFunc("GET /event/{id}", handler.GetById)
	r.HandleFunc("POST /event/", handler.Create)

	r.HandleFunc("PUT /event/{id}", middlewares.AdminOnly(handler.Update))
	r.HandleFunc("PATCH /event/{id}", middlewares.AdminOnly(handler.Update))
	r.HandleFunc("PATCH /event/{id}/approve", middlewares.AdminOnly(handler.Approve))

	r.HandleFunc("DELETE /event/{id}", handler.DeleteById)
}
