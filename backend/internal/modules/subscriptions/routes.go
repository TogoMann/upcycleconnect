package subscriptions

import (
	"backend/internal/middlewares"
	"backend/internal/modules/financial"
	"backend/internal/modules/plans"
	"backend/internal/modules/users"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	userRepo := users.NewRepository(db)
	planRepo := plans.NewRepository(db)
	finRepo := financial.NewRepository(db)
	finSvc := financial.NewService(finRepo)
	repo := NewRepository(db)
	service := NewService(repo, userRepo, planRepo, finSvc)
	handler := NewHandler(service)

	r.HandleFunc("GET /admin/abonnements", middlewares.AdminOnly(handler.GetAllAbonnements))
	r.HandleFunc("GET /subscriptions", middlewares.AdminOnly(handler.GetAll))
	r.HandleFunc("GET /subscriptions/me", middlewares.Authenticated(handler.GetMySubscription))
	r.HandleFunc("GET /subscriptions/{id}", middlewares.Authenticated(handler.GetById))
	r.HandleFunc("POST /subscriptions/", middlewares.AdminOnly(handler.Create))
	r.HandleFunc("POST /subscriptions/choose", middlewares.Authenticated(handler.ChoosePlan))
	r.HandleFunc("PUT /subscriptions/{id}", middlewares.ProOnly(handler.Update))
	r.HandleFunc("PATCH /subscriptions/{id}", middlewares.ProOnly(handler.Update))
	r.HandleFunc("DELETE /subscriptions/{id}", middlewares.AdminOnly(handler.DeleteById))
}
