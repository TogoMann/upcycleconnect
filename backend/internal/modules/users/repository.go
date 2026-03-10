package tables

import (
	"API/db"
	"API/models"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func GetAllUsers() ([]models.User, error) {
	rows, err := db.Conn.Query(db.Ctx, "SELECT id, first_name, last_name, email, password_hash, role FROM users")
	if err != nil {
		return nil, fmt.Errorf("package tables/users GetAllusers query: %w", err)
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.User])

	if err != nil {
		return nil, fmt.Errorf("package tables/users GetAllUsers: %v", err.Error())
	}

	return users, nil
}

func GetUserById(id string) (*models.User, error) {
	rows, err := db.Conn.Query(db.Ctx, "SELECT id, first_name, last_name, email, password_hash, role FROM users WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("package tables/users GetUserById query: %w", err)
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.User])

	if err != nil {
		return nil, fmt.Errorf("package tables/users GetUserById: %v", err.Error())
	}
	return &user, nil
}
