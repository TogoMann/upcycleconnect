package payments

import (
	"backend/internal/middlewares"
	"backend/internal/modules/advertisement"
	"backend/internal/modules/listing"
	listingorder "backend/internal/modules/listing_order"
	"backend/internal/modules/plans"
	"backend/internal/modules/subscriptions"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool, listingService *listing.Service, listingOrderService *listingorder.Service, advertisementService *advertisement.Service, planService *plans.Service, subscriptionService *subscriptions.Service) {
	service := NewService(listingService, listingOrderService, advertisementService, planService, subscriptionService)
	handler := NewHandler(service)

	r.HandleFunc("POST /subscriptions/checkout", middlewares.Authenticated(handler.CreateSubscriptionCheckout))
	r.HandleFunc("POST /advertisement/{id}/checkout", middlewares.Authenticated(handler.CreateAdvertisementCheckout))
	r.HandleFunc("POST /listing-order/checkout", middlewares.Authenticated(handler.CreateListingOrderCheckout))
	r.HandleFunc("GET /payments/verify", middlewares.Authenticated(handler.VerifySession))
	r.HandleFunc("POST /webhooks/stripe", handler.Webhook)
}
