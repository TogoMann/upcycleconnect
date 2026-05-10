package cart

import (
	"backend/internal/middlewares"
	listingorder "backend/internal/modules/listing_order"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	orderRepo := listingorder.NewRepository(db)
	orderService := listingorder.NewService(orderRepo)

	repo := NewRepository(db)
	service := NewService(repo, orderService)
	handler := NewHandler(service)

	r.HandleFunc("GET /cart", middlewares.Authenticated(handler.Get))
	r.HandleFunc("POST /cart", middlewares.Authenticated(handler.Add))
	r.HandleFunc("DELETE /cart/{listingId}", middlewares.Authenticated(handler.Remove))
	r.HandleFunc("POST /cart/checkout", middlewares.Authenticated(handler.Checkout))
}
