package entryparticipation

import (
	"net/http"
	
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /entry_participation", handler.GetAll)
	r.HandleFunc("GET /entry_participation/{id}", handler.GetById)
	r.HandleFunc("POST /entry_participation/", handler.Create)
	r.HandleFunc("DELETE /entry_participation/{id}", handler.DeleteById)
}
