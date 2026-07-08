package listing

import (
	"backend/internal/middlewares"
	"backend/internal/modules/container"
	"backend/internal/modules/item"
	"backend/internal/modules/subscriptions"
	"backend/internal/modules/users"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {

	repo := NewRepository(db)
	subRepo := subscriptions.NewRepository(db)

	userRepo := users.NewRepository(db)
	userService := users.NewService(userRepo)

	containerRepo := container.NewRepository(db)
	containerService := container.NewService(containerRepo)

	itemRepo := item.NewRepository(db)
	itemService := item.NewService(itemRepo, userService)

	service := NewService(repo, subRepo, containerService, itemService, userService)

	handler := NewHandler(service, userService)

	r.HandleFunc("GET /listing", middlewares.OptionalAuth(handler.GetAllApproved))
	r.HandleFunc("GET /listing/me", middlewares.Authenticated(handler.GetByUserId))
	r.HandleFunc("GET /admin/listings", middlewares.AdminOnly(handler.GetAll))
	r.HandleFunc("GET /admin/listings/{id}", middlewares.AdminOnly(handler.GetById))
	r.HandleFunc("GET /listing/{id}", handler.GetById)

	r.HandleFunc("POST /listing/", middlewares.Authenticated(handler.Create))
	r.HandleFunc("POST /listing/upload", middlewares.Authenticated(handler.UploadImage))

	r.HandleFunc("PUT /listing/{id}", middlewares.AdminOnly(handler.Update))
	r.HandleFunc("PATCH /listing/{id}", middlewares.AdminOnly(handler.Update))
	r.HandleFunc("PATCH /listing/{id}/approve", middlewares.AdminOnly(handler.Approve))
	r.HandleFunc("PATCH /listing/{id}/disapprove", middlewares.AdminOnly(handler.Disapprove))

	r.HandleFunc("DELETE /listing/{id}", middlewares.Authenticated(handler.DeleteById))
}
