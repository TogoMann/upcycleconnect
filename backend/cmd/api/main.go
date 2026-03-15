package main

import (
	db "backend/internal/database"
	"backend/internal/router"
	"fmt"
	"net/http"
)

func main() {
	db.Conn = db.NewDB()

	r := router.NewRouter(db.Conn)

	//	http.HandleFunc("GET /users/{$}", app.GetAllUsers)
	//http.HandleFunc("GET /users/{id}", app.GetUserById)
	//	http.HandleFunc("POST /users/{$}", app.CreateUser)
	//	http.HandleFunc("POST /login/{$}", app.Login)

	fmt.Println("Listening at http://localhost:8081")
	http.ListenAndServe(":8081", r)
}
