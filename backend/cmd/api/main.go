package main

import (
	db "backend/internal/database"
	"backend/internal/middlewares"
	"backend/internal/router"
	"fmt"
	"net/http"
)

func main() {
	db.Conn = db.NewDB()

	r := router.NewRouter(db.Conn)
	handlerWithCors := middlewares.CorsMiddleware(r)

	fmt.Println("Listening at http://localhost:8081")
	http.ListenAndServe(":8081", handlerWithCors)
}
