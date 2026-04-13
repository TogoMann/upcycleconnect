package course

import (
	"backend/internal/middlewares"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /course", handler.GetAll)
	r.HandleFunc("GET /course/{id}", handler.GetById)

	r.HandleFunc("GET /users/{user_id}/courses", handler.GetUserCourses)

	r.HandleFunc("POST /course/", handler.Create)

	r.HandleFunc("PUT /course/{id}", middlewares.AdminOnly(handler.Update))
	r.HandleFunc("PATCH /course/{id}", middlewares.AdminOnly(handler.Update))

	r.HandleFunc("DELETE /course/{id}", handler.DeleteById)
}
