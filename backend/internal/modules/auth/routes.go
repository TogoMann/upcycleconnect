package auth

import (
	"backend/internal/modules/users"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	userRepo := users.NewRepository(db)
	service := NewService(userRepo)
	handler := NewHandler(service)

	r.HandleFunc("POST /login", handler.Login)
}
