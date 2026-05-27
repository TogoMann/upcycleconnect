package tests

import (
	db "backend/internal/database"
	"backend/internal/router"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// SetupTestRouter initializes the router with the real DB pool (for integration tests)
// In a real scenario, we might want to use a separate test database.
func SetupTestRouter(pool *pgxpool.Pool) http.Handler {
	return router.NewRouter(pool)
}

func GetPool() *pgxpool.Pool {
	if db.Pool == nil {
		// Load .env if not already loaded (e.g. running from run_tests.sh on host)
		if os.Getenv("DB_USER") == "" {
			godotenv.Load("../.env") // Try to load from backend/.env if it exists
			godotenv.Load("../../.env") // Or root if running from backend/tests/
		}

		// Override DB_HOST for local tests (outside Docker)
		os.Setenv("DB_HOST", "localhost")
		db.Pool = db.NewDB()
	}
	return db.Pool
}
