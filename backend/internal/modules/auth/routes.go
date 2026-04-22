package auth

import (
	"backend/internal/modules/subscriptions"
	"backend/internal/modules/users"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	userRepo := users.NewRepository(db)
	subRepo := subscriptions.NewRepository(db)
	service := NewService(userRepo, subRepo)
	handler := NewHandler(service)

	r.HandleFunc("POST /login/", handler.Login)
	r.HandleFunc("POST /register/", handler.Register)
}
