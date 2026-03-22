package eventparticipation

import (
	"net/http"

	
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /event-participation", handler.GetAll)
	r.HandleFunc("GET /event-participation/{id}", handler.GetById)
	r.HandleFunc("POST /event-participation/", handler.Create)
	r.HandleFunc("DELETE /event-participation/{id}", handler.DeleteById)
}
