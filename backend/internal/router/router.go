package router

import (
	db "backend/internal/database"
	"backend/internal/modules/post"
	"backend/internal/modules/thread"
	"backend/internal/modules/users"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	err := db.Conn.Ping(db.Ctx)
	if err != nil {
		http.Error(w, " ping fail", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	res, _ := json.Marshal("en vie.")
	fmt.Fprintf(w, "%s", string(res))
}

func NewRouter(db *pgx.Conn) *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("GET /", healthCheck)

	// modules
	users.RegisterRoutes(r, db)
	thread.RegisterRoutes(r, db)
	post.RegisterRoutes(r, db)

	return r
}
