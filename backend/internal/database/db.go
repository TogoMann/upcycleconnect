package db

import (
	"backend/internal/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool
var Ctx = context.Background()

func NewDB() *pgxpool.Pool {
	cfg := config.Load()
	cred := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
	conn, err := pgxpool.New(Ctx, cred)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Database connected")
	return conn
}
