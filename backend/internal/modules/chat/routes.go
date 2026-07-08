package chat

import (
	"backend/internal/middlewares"
	"backend/internal/modules/container"
	"backend/internal/modules/course"
	"backend/internal/modules/item"
	"backend/internal/modules/listing"
	"backend/internal/modules/subscriptions"
	"backend/internal/modules/users"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	userRepo := users.NewRepository(db)
	userService := users.NewService(userRepo)

	containerRepo := container.NewRepository(db)
	containerService := container.NewService(containerRepo)
	itemRepo := item.NewRepository(db)
	itemService := item.NewService(itemRepo, userService)

	listingRepo := listing.NewRepository(db)
	subRepo := subscriptions.NewRepository(db)
	listingService := listing.NewService(listingRepo, subRepo, containerService, itemService, userService)

	courseRepo := course.NewRepository(db)
	courseService := course.NewService(courseRepo)

	repo := NewRepository(db)
	service := NewService(repo, listingService, courseService)
	handler := NewHandler(service)

	r.HandleFunc("POST /chat/messages", middlewares.Authenticated(handler.SendMessage))
	r.HandleFunc("GET /chat/conversations", middlewares.Authenticated(handler.GetUserConversations))
	r.HandleFunc("GET /chat/conversations/{id}/messages", middlewares.Authenticated(handler.GetConversationMessages))
	r.HandleFunc("PUT /chat/messages/{id}", middlewares.Authenticated(handler.EditMessage))
	r.HandleFunc("POST /chat/messages/{id}/proposal", middlewares.Authenticated(handler.HandleProposal))

	r.HandleFunc("GET /admin/chat/conversations", middlewares.AdminOnly(handler.AdminGetConversations))
	r.HandleFunc("GET /admin/chat/conversations/{id}", middlewares.AdminOnly(handler.AdminGetConversationDetails))
	r.HandleFunc("DELETE /admin/chat/messages/{id}", middlewares.AdminOnly(handler.AdminDeleteMessage))
}
