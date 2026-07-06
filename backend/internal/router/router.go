package router

import (
	db "backend/internal/database"
	"backend/internal/modules/advertisement"
	"backend/internal/modules/auth"
	"backend/internal/modules/brands"
	"backend/internal/modules/cart"
	"backend/internal/modules/chat"
	"backend/internal/modules/city"
	"backend/internal/modules/comments"
	"backend/internal/modules/companies"
	"backend/internal/modules/container"
	"backend/internal/modules/course"
	courseorder "backend/internal/modules/course_order"
	"backend/internal/modules/entry"
	entryparticipation "backend/internal/modules/entry_participation"
	"backend/internal/modules/event"
	eventparticipation "backend/internal/modules/event_participation"
	"backend/internal/modules/financial"
	"backend/internal/modules/item"
	"backend/internal/modules/listing"
	listingorder "backend/internal/modules/listing_order"
	"backend/internal/modules/logs"
	"backend/internal/modules/news"
	"backend/internal/modules/notifications"
	paymentmethods "backend/internal/modules/payment_methods"
	"backend/internal/modules/payments"
	"backend/internal/modules/planning"
	"backend/internal/modules/plans"
	"backend/internal/modules/post"
	"backend/internal/modules/project"
	"backend/internal/modules/reporting"
	"backend/internal/modules/settings"
	"backend/internal/modules/stats"
	"backend/internal/modules/subscriptions"
	"backend/internal/modules/thread"
	"backend/internal/modules/users"

	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	err := db.Pool.Ping(db.Ctx)
	if err != nil {
		http.Error(w, " ping fail", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("en vie.")
}

func NewRouter(db *pgxpool.Pool) *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("GET /health", healthCheck)

	fs := http.FileServer(http.Dir("./uploads"))
	r.Handle("GET /uploads/", http.StripPrefix("/uploads/", fs))

	auth.RegisterRoutes(r, db)
	city.RegisterRoutes(r, db)
	companies.RegisterRoutes(r, db)
	advertisement.RegisterRoutes(r, db)
	users.RegisterRoutes(r, db)
	item.RegisterRoutes(r, db)
	thread.RegisterRoutes(r, db)
	post.RegisterRoutes(r, db)
	news.RegisterRoutes(r, db)
	listingorder.RegisterRoutes(r, db)
	listing.RegisterRoutes(r, db)
	cart.RegisterRoutes(r, db)
	notifications.RegisterRoutes(r, db)
	planning.RegisterRoutes(r, db)
	plans.RegisterRoutes(r, db)
	eventparticipation.RegisterRoutes(r, db)
	event.RegisterRoutes(r, db)
	entryparticipation.RegisterRoutes(r, db)
	entry.RegisterRoutes(r, db)
	courseorder.RegisterRoutes(r, db)
	course.RegisterRoutes(r, db)
	subscriptions.RegisterRoutes(r, db)
	comments.RegisterRoutes(r, db)
	container.RegisterRoutes(r, db)
	financial.RegisterRoutes(r, db)
	project.RegisterRoutes(r, db)
	chat.RegisterRoutes(r, db)
	reporting.RegisterRoutes(r, db)
	brands.RegisterRoutes(r, db)
	paymentmethods.RegisterRoutes(r, db)
	logs.RegisterRoutes(r, db)

	userRepo := users.NewRepository(db)
	userSvc := users.NewService(userRepo)

	containerRepo := container.NewRepository(db)
	containerSvc := container.NewService(containerRepo)
	itemRepo := item.NewRepository(db)
	itemSvc := item.NewService(itemRepo, userSvc)

	listingRepo := listing.NewRepository(db)
	subRepo := subscriptions.NewRepository(db)
	listingSvc := listing.NewService(listingRepo, subRepo, containerSvc, itemSvc, userSvc)

	chatRepo := chat.NewRepository(db)
	finRepo := financial.NewRepository(db)
	finSvc := financial.NewService(finRepo)
	listingOrderRepo := listingorder.NewRepository(db)
	listingOrderSvc := listingorder.NewService(listingOrderRepo, finSvc, listingSvc, chatRepo, containerSvc, itemSvc, userSvc)

	advertisementRepo := advertisement.NewRepository(db)
	advertisementSvc := advertisement.NewService(advertisementRepo)

	planRepo := plans.NewRepository(db)
	planSvc := plans.NewService(planRepo)

	subscriptionRepo := subscriptions.NewRepository(db)
	subscriptionSvc := subscriptions.NewService(subscriptionRepo, userRepo, planRepo, finSvc)

	payments.RegisterRoutes(r, db, listingSvc, listingOrderSvc, advertisementSvc, planSvc, subscriptionSvc)
	stats.RegisterRoutes(r, db)
	settings.RegisterRoutes(r, db)

	return r
}
