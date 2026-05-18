package project

import (
	"backend/internal/middlewares"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.HandleFunc("GET /admin/projets", middlewares.AdminOnly(handler.GetAllProjets))
	r.HandleFunc("PATCH /admin/projets/{id}/mise-en-avant", middlewares.AdminOnly(handler.UpdateFeatured))

	r.HandleFunc("GET /project", handler.GetAll)
	r.HandleFunc("GET /project/{id}", handler.GetById)
	r.HandleFunc("GET /project/{id}/steps", handler.GetSteps)

	r.HandleFunc("POST /project", middlewares.HasPlan("Pro")(handler.Create))
	r.HandleFunc("PUT /project/{id}", middlewares.HasPlan("Pro")(handler.Update))
	r.HandleFunc("PATCH /project/{id}", middlewares.HasPlan("Pro")(handler.Update))
	r.HandleFunc("DELETE /project/{id}", middlewares.HasPlan("Pro")(handler.DeleteById))

	r.HandleFunc("POST /project/steps", middlewares.HasPlan("Pro")(handler.CreateStep))
	r.HandleFunc("PUT /project/steps/{step_id}", middlewares.HasPlan("Pro")(handler.UpdateStep))
	r.HandleFunc("DELETE /project/steps/{step_id}", middlewares.HasPlan("Pro")(handler.DeleteStep))
}
