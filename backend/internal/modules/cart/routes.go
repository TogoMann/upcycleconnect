package cart

import (
	"backend/internal/middlewares"
	"backend/internal/modules/chat"
	"backend/internal/modules/financial"
	"backend/internal/modules/listing"
	listingorder "backend/internal/modules/listing_order"
	"backend/internal/modules/subscriptions"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	finRepo := financial.NewRepository(db)
	finSvc := financial.NewService(finRepo)

	subRepo := subscriptions.NewRepository(db)
	listingRepo := listing.NewRepository(db)
	listingSvc := listing.NewService(listingRepo, subRepo)

	chatRepo := chat.NewRepository(db)

	orderRepo := listingorder.NewRepository(db)
	orderService := listingorder.NewService(orderRepo, finSvc, listingSvc, chatRepo)

	repo := NewRepository(db)
	service := NewService(repo, orderService)
	handler := NewHandler(service)

	r.HandleFunc("GET /cart", middlewares.Authenticated(handler.Get))
	r.HandleFunc("POST /cart", middlewares.Authenticated(handler.Add))
	r.HandleFunc("DELETE /cart/{listingId}", middlewares.Authenticated(handler.Remove))
	r.HandleFunc("POST /cart/checkout", middlewares.Authenticated(handler.Checkout))
}
