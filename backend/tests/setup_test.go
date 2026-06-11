package tests

import (
	db "backend/internal/database"
	"backend/internal/router"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func SetupTestRouter(pool *pgxpool.Pool) http.Handler {
	return router.NewRouter(pool)
}

func GetPool() *pgxpool.Pool {
	if db.Pool == nil {
		if os.Getenv("DB_USER") == "" {
			godotenv.Load("../.env")
			godotenv.Load("../../.env")
		}

		os.Setenv("DB_HOST", "localhost")
		db.Pool = db.NewDB()
	}
	return db.Pool
}
