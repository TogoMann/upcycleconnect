package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn
var Ctx = context.Background()

func NewDB() *pgx.Conn {
	conn, err := pgx.Connect(Ctx, "postgres://admin:root@localhost:5432/upcycleconnect")

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Database connected")
	return conn
}
