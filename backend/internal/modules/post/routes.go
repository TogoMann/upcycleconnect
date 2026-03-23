package post

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /post", handler.GetAll)
	r.HandleFunc("GET /post/{id}", handler.GetById)

	// Get all posts in a thread
	r.HandleFunc("GET /thread/{thread_id}/posts", handler.GetThreadPosts)

	r.HandleFunc("POST /post/", handler.Create)
	r.HandleFunc("PUT /post/{id}", handler.Update)

	r.HandleFunc("DELETE /post/{id}", handler.DeleteById)
}
