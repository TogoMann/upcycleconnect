package post

import (
	"backend/internal/middlewares"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /post", handler.GetAll)
	r.HandleFunc("GET /post/{id}", handler.GetById)
	r.HandleFunc("GET /thread/{thread_id}/posts", handler.GetThreadPosts)
	r.HandleFunc("POST /post", handler.Create)
	r.HandleFunc("PUT /post/{id}", middlewares.StaffOnly(handler.Update))
	r.HandleFunc("PATCH /post/{id}", middlewares.StaffOnly(handler.Update))
	r.HandleFunc("DELETE /post/{id}", middlewares.StaffOnly(handler.DeleteById))
}
