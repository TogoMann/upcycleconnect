package thread

import (
	"backend/internal/middlewares"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /thread", handler.GetAll)
	r.HandleFunc("GET /thread/{id}", handler.GetById)
	r.HandleFunc("POST /thread", middlewares.Authenticated(handler.Create))
	r.HandleFunc("POST /thread/{id}/upvote", middlewares.Authenticated(handler.Upvote))
	r.HandleFunc("POST /thread/{id}/downvote", middlewares.Authenticated(handler.Downvote))
	r.HandleFunc("PUT /thread/{id}", middlewares.StaffOnly(handler.Update))
	r.HandleFunc("PATCH /thread/{id}", middlewares.StaffOnly(handler.Update))
	r.HandleFunc("DELETE /thread/{id}", middlewares.StaffOnly(handler.DeleteById))

	r.HandleFunc("GET /salarie/forum", middlewares.StaffOnly(handler.GetSalarieForum))
	r.HandleFunc("POST /salarie/forum/{id}/epingler", middlewares.StaffOnly(handler.Epingler))
	r.HandleFunc("POST /salarie/forum/bannir", middlewares.StaffOnly(handler.Bannir))
	r.HandleFunc("DELETE /salarie/forum/{id}", middlewares.StaffOnly(handler.DeleteById))
}
