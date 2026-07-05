package auth

import (
	"backend/internal/middlewares"
	"backend/internal/modules/companies"
	"backend/internal/modules/plans"
	"backend/internal/modules/subscriptions"
	"backend/internal/modules/users"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	userRepo := users.NewRepository(db)
	userService := users.NewService(userRepo)
	subRepo := subscriptions.NewRepository(db)
	planRepo := plans.NewRepository(db)
	compRepo := companies.NewRepository(db)
	service := NewService(userRepo, userService, subRepo, planRepo, compRepo)
	handler := NewHandler(service)

	r.HandleFunc("POST /login/", handler.Login)
	r.HandleFunc("POST /login", handler.Login)
	r.HandleFunc("POST /register/", handler.Register)
	r.HandleFunc("POST /register", handler.Register)
	r.HandleFunc("POST /auth/admin/reset-request", middlewares.AdminOnly(handler.AdminRequestReset))
	r.HandleFunc("POST /auth/forgot-password", handler.ForgotPassword)
	r.HandleFunc("POST /auth/reset-password", handler.ResetPassword)
	r.HandleFunc("GET /siret/verify", handler.VerifySiret)
}
