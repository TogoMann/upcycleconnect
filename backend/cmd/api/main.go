package main

import (
	"API/app"
	"API/db"
	"encoding/json"
	"fmt"
	"net/http"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	err := db.Conn.Ping(db.Ctx)
	if err != nil {
		http.Error(w, " ping fail", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("school", "esgi")

	res, _ := json.Marshal("en vie.")
	fmt.Fprintf(w, "%s", string(res))
}

func main() {
	db.Conn = db.NewDB()

	http.HandleFunc("GET /{$}", healthCheck)
	http.HandleFunc("GET /users/{$}", app.GetAllUsers)
	http.HandleFunc("GET /users/{id}", app.GetUserById)
	//	http.HandleFunc("POST /users/{$}", app.CreateUser)
	//	http.HandleFunc("POST /login/{$}", app.Login)

	fmt.Println("Listening at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
