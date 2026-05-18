package main

import (
	"backend/internal/config"
	db "backend/internal/database"
	"backend/internal/middlewares"
	"backend/internal/router"
	"fmt"
	"net/http"
)

func main() {
	cfg := config.Load()
	db.Pool = db.NewDB()

	r := router.NewRouter(db.Pool)
	handlerWithCors := middlewares.CorsMiddleware(r)

	port := cfg.AppPort
	if port == "" {
		port = "8081"
	}

	fmt.Printf("Listening at http://localhost:%s\n", port)
	http.ListenAndServe(":"+port, handlerWithCors)
}
