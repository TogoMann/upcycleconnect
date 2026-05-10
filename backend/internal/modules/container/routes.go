package container

import (
	"backend/internal/middlewares"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /admin/conteneurs", middlewares.AdminOnly(handler.GetAll))
	r.HandleFunc("GET /container/{id}", middlewares.AdminOnly(handler.GetById))
	r.HandleFunc("GET /container/{id}/lockers", handler.GetLockers)
	r.HandleFunc("POST /admin/lockers", middlewares.AdminOnly(handler.CreateLocker))
	r.HandleFunc("PUT /admin/lockers/{id}", middlewares.AdminOnly(handler.UpdateLocker))
	r.HandleFunc("DELETE /admin/lockers/{id}", middlewares.AdminOnly(handler.DeleteLocker))
	r.HandleFunc("POST /container", middlewares.AdminOnly(handler.Create))
	r.HandleFunc("PUT /container/{id}", middlewares.AdminOnly(handler.Update))
	r.HandleFunc("DELETE /container/{id}", middlewares.AdminOnly(handler.Delete))
}
