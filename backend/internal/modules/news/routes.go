package news

import (
	"backend/internal/middlewares"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /news", handler.GetAllPublic)
	r.HandleFunc("GET /salarie/conseils", handler.GetAll)
	r.HandleFunc("GET /news/{id}", handler.GetById)
	r.HandleFunc("POST /news", middlewares.StaffOnly(handler.Create))
	r.HandleFunc("POST /salarie/conseils", middlewares.StaffOnly(handler.Create))
	r.HandleFunc("DELETE /news/{id}", middlewares.StaffOnly(handler.DeleteById))
	r.HandleFunc("DELETE /salarie/conseils/{id}", middlewares.StaffOnly(handler.DeleteById))
}
