package db

import (
	"backend/internal/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn
var Ctx = context.Background()

func NewDB() *pgx.Conn {
	cfg := config.Load()
	cred := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
	conn, err := pgx.Connect(Ctx, cred)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Database connected")
	return conn
}
