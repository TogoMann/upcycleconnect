package router

import (
	db "backend/internal/database"
	"backend/internal/modules/advertisement"
	"backend/internal/modules/auth"
	"backend/internal/modules/city"
	"backend/internal/modules/comments"
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
	"backend/internal/modules/news"
	"backend/internal/modules/planning"
	"backend/internal/modules/plans"
	"backend/internal/modules/post"
	"backend/internal/modules/project"
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
	advertisement.RegisterRoutes(r, db)
	users.RegisterRoutes(r, db)
	item.RegisterRoutes(r, db)
	thread.RegisterRoutes(r, db)
	post.RegisterRoutes(r, db)
	news.RegisterRoutes(r, db)
	listingorder.RegisterRoutes(r, db)
	listing.RegisterRoutes(r, db)
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

	return r
}
