package listingorder

import (
	"backend/internal/middlewares"
	"backend/internal/modules/chat"
	"backend/internal/modules/container"
	"backend/internal/modules/financial"
	"backend/internal/modules/item"
	"backend/internal/modules/listing"
	"backend/internal/modules/subscriptions"
	"backend/internal/modules/users"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {

	repo := NewRepository(db)
	finRepo := financial.NewRepository(db)
	finSvc := financial.NewService(finRepo)

	chatRepo := chat.NewRepository(db)

	userRepo := users.NewRepository(db)
	userSvc := users.NewService(userRepo)

	containerRepo := container.NewRepository(db)
	containerSvc := container.NewService(containerRepo)

	itemRepo := item.NewRepository(db)
	itemSvc := item.NewService(itemRepo, userSvc)

	subRepo := subscriptions.NewRepository(db)
	listingRepo := listing.NewRepository(db)
	listingSvc := listing.NewService(listingRepo, subRepo, containerSvc, itemSvc, userSvc)

	service := NewService(repo, finSvc, listingSvc, chatRepo, containerSvc, itemSvc, userSvc)
	handler := NewHandler(service)

	r.HandleFunc("GET /listing-order/me", middlewares.Authenticated(handler.GetByUserId))
	r.HandleFunc("GET /listing-order", middlewares.AdminOnly(handler.GetAll))
	r.HandleFunc("GET /listing-order/{id}", middlewares.Authenticated(handler.GetById))

	r.HandleFunc("POST /listing-order/", middlewares.Authenticated(handler.Create))

	r.HandleFunc("DELETE /listing-order/{id}", middlewares.AdminOnly(handler.DeleteById))
}
