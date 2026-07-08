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

	r.HandleFunc("GET /news", middlewares.OptionalAuth(handler.GetAllActualites))
	r.HandleFunc("GET /conseils", middlewares.OptionalAuth(handler.GetAllConseils))
	r.HandleFunc("GET /salarie/conseils", middlewares.OptionalAuth(handler.GetAllConseils))
	r.HandleFunc("GET /news/{id}", handler.GetById)
	r.HandleFunc("POST /news", middlewares.StaffOnly(handler.CreateActualite))
	r.HandleFunc("POST /salarie/conseils", middlewares.StaffOnly(handler.CreateConseil))
	r.HandleFunc("PUT /news/{id}", middlewares.StaffOnly(handler.Update))
	r.HandleFunc("PUT /salarie/conseils/{id}", middlewares.StaffOnly(handler.Update))
	r.HandleFunc("DELETE /news/{id}", middlewares.StaffOnly(handler.DeleteById))
	r.HandleFunc("DELETE /salarie/conseils/{id}", middlewares.StaffOnly(handler.DeleteById))
	r.HandleFunc("POST /news/{id}/vote", middlewares.Authenticated(handler.Vote))
	r.HandleFunc("POST /conseils/{id}/vote", middlewares.Authenticated(handler.Vote))
}
