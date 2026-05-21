package courseorder

import (
	"backend/internal/middlewares"
	"backend/internal/modules/course"
	"backend/internal/modules/financial"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *http.ServeMux, db *pgxpool.Pool) {

	repo := NewRepository(db)
	finRepo := financial.NewRepository(db)
	finSvc := financial.NewService(finRepo)

	courseRepo := course.NewRepository(db)
	courseSvc := course.NewService(courseRepo)

	service := NewService(repo, finSvc, courseSvc)
	handler := NewHandler(service)

	r.HandleFunc("GET /course-order/me", middlewares.Authenticated(handler.GetByUserId))
	r.HandleFunc("GET /course-order", middlewares.AdminOnly(handler.GetAll))
	r.HandleFunc("GET /course-order/{id}", middlewares.Authenticated(handler.GetById))

	r.HandleFunc("POST /course-order/", middlewares.Authenticated(handler.Create))

	r.HandleFunc("DELETE /course-order/{id}", middlewares.AdminOnly(handler.DeleteById))
}
