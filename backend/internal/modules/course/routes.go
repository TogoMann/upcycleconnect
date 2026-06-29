package course

import (
	"backend/internal/middlewares"
	"backend/internal/modules/users"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := NewRepository(db)
	service := NewService(repo)

	userRepo := users.NewRepository(db)
	userService := users.NewService(userRepo)

	handler := NewHandler(service, userService)

	r.HandleFunc("GET /course/catalogue", handler.GetAllApprovedCatalogue)
	r.HandleFunc("GET /admin/catalogue", middlewares.AdminOnly(handler.GetAllCatalogue))
	r.HandleFunc("POST /admin/catalogue", middlewares.AdminOnly(handler.Create))
	r.HandleFunc("PUT /admin/catalogue/{id}", middlewares.AdminOnly(handler.Update))
	r.HandleFunc("PATCH /admin/catalogue/{id}/approve", middlewares.AdminOnly(handler.Approve))
	r.HandleFunc("PATCH /admin/catalogue/{id}/propose", middlewares.AdminOnly(handler.ProposeModification))
	r.HandleFunc("PATCH /admin/catalogue/{id}/disapprove", middlewares.AdminOnly(handler.Disapprove))
	r.HandleFunc("DELETE /admin/catalogue/{id}", middlewares.AdminOnly(handler.DeleteById))

	r.HandleFunc("GET /course", handler.GetAll)
	r.HandleFunc("GET /course/{id}", handler.GetById)
	r.HandleFunc("GET /users/{user_id}/courses", handler.GetUserCourses)
	r.HandleFunc("POST /course", middlewares.Authenticated(handler.Create))
	r.HandleFunc("PUT /course/{id}", middlewares.Authenticated(handler.Update))
	r.HandleFunc("PATCH /course/{id}", middlewares.Authenticated(handler.Update))
	r.HandleFunc("PATCH /course/{id}/approve", middlewares.AdminOnly(handler.Approve))
	r.HandleFunc("DELETE /course/{id}", middlewares.Authenticated(handler.DeleteById))

	r.HandleFunc("GET /salarie/formations", middlewares.Authenticated(handler.GetMyCourses))
	r.HandleFunc("GET /salarie/formations/{id}", middlewares.Authenticated(handler.GetById))
	r.HandleFunc("POST /salarie/formations", middlewares.Authenticated(handler.Create))
	r.HandleFunc("PUT /salarie/formations/{id}", middlewares.Authenticated(handler.Update))
	r.HandleFunc("DELETE /salarie/formations/{id}", middlewares.Authenticated(handler.DeleteById))
}
