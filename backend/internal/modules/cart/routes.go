package cart

import (
	"backend/internal/middlewares"
	"backend/internal/modules/chat"
	"backend/internal/modules/container"
	"backend/internal/modules/course"
	courseorder "backend/internal/modules/course_order"
	eventparticipation "backend/internal/modules/event_participation"
	"backend/internal/modules/financial"
	"backend/internal/modules/item"
	"backend/internal/modules/listing"
	listingorder "backend/internal/modules/listing_order"
	"backend/internal/modules/subscriptions"
	"backend/internal/modules/users"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	finRepo := financial.NewRepository(db)
	finSvc := financial.NewService(finRepo)

	userRepo := users.NewRepository(db)
	userSvc := users.NewService(userRepo)

	containerRepo := container.NewRepository(db)
	containerSvc := container.NewService(containerRepo)

	itemRepo := item.NewRepository(db)
	itemSvc := item.NewService(itemRepo, userSvc)

	subRepo := subscriptions.NewRepository(db)
	listingRepo := listing.NewRepository(db)
	listingSvc := listing.NewService(listingRepo, subRepo, containerSvc, itemSvc, userSvc)

	chatRepo := chat.NewRepository(db)

	orderRepo := listingorder.NewRepository(db)
	orderService := listingorder.NewService(orderRepo, finSvc, listingSvc, chatRepo, containerSvc, itemSvc, userSvc)

	eventPartRepo := eventparticipation.NewRepository(db)

	courseRepo := course.NewRepository(db)
	courseSvc := course.NewService(courseRepo)
	courseOrdRepo := courseorder.NewRepository(db)
	courseOrdSvc := courseorder.NewService(courseOrdRepo, finSvc, courseSvc, userSvc)

	repo := NewRepository(db)
	service := NewService(repo, orderService, eventPartRepo, courseOrdSvc, finSvc)
	handler := NewHandler(service)

	r.HandleFunc("GET /cart", middlewares.Authenticated(handler.Get))
	r.HandleFunc("POST /cart", middlewares.Authenticated(handler.Add))
	r.HandleFunc("DELETE /cart/{type}/{id}", middlewares.Authenticated(handler.Remove))
	r.HandleFunc("DELETE /cart", middlewares.Authenticated(handler.Clear))
	r.HandleFunc("POST /cart/checkout", middlewares.Authenticated(handler.Checkout))
}
