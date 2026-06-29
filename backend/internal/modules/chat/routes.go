package chat

import (
	"backend/internal/middlewares"
	"backend/internal/modules/listing"
	"backend/internal/modules/subscriptions"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	listingRepo := listing.NewRepository(db)
	subRepo := subscriptions.NewRepository(db)
	listingService := listing.NewService(listingRepo, subRepo)

	repo := NewRepository(db)
	service := NewService(repo, listingService)
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
