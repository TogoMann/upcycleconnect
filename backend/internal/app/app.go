package app

import (
	"API/db/tables"
	"API/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := tables.GetAllUsers()

	if err != nil {
		fmt.Println(err.Error())
	}

	res, _ := json.Marshal(users)
	fmt.Fprintf(w, "%s", string(res))
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("id")

	w.Header().Set("Content-Type", "application/json")

	user, err := tables.GetUserById(userId)

	if err != nil {
		fmt.Println("Erreur DB GetUserById:", err.Error())
		http.Error(w, "Couldn't get user", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, fmt.Sprintf("No User with id: %s", userId), http.StatusNotFound)
		return
	}

	res, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Erreur d'encodage JSON du résultat", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", string(res))
}

func ValidateUser(userDto models.User) []string {
	errsMsg := []string{}

	if len(userDto.FirstName) < 1 || len(userDto.LastName) < 1 {
		errsMsg = append(errsMsg, "Longueur du pseudonyme doit être de 5 caractères minimum.")
	}

	// if len(userDto.Password) < 8 {
	// 	errsMsg = append(errsMsg, "Longueur du mot de passe doit être de 8 caractères minimum.")
	// }

	// if !strings.ContainsAny(userDto.Password, "!-$+/") {
	// 	errsMsg = append(errsMsg, "Le mot de passe doit contenir au moins 1 caractère spécial.")
	// }

	return errsMsg
}
