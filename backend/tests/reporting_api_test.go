package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestReportingAPI(t *testing.T) {
	pool := GetPool()
	router := SetupTestRouter(pool)
	adminToken := getAdminTokenForReportingAPI(t, router, pool)

	t.Run("GetActorStats", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/reporting/actors", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d: %s", rr.Code, rr.Body.String())
		}

		var stats []map[string]interface{}
		if err := json.Unmarshal(rr.Body.Bytes(), &stats); err != nil {
			t.Fatal(err)
		}
		if len(stats) == 0 {
			t.Error("expected actor stats, got empty list")
		}
	})

	t.Run("GetPrestationStats", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/reporting/prestations", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d: %s", rr.Code, rr.Body.String())
		}

		var stats []map[string]interface{}
		if err := json.Unmarshal(rr.Body.Bytes(), &stats); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("GetUserPredictions", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/reporting/predictions", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d: %s", rr.Code, rr.Body.String())
		}
	})
}

func getAdminTokenForReportingAPI(t *testing.T, router http.Handler, pool *pgxpool.Pool) string {
	adminUsername := "admin_reporting_test"
	adminEmail := "admin_reporting_test@example.com"
	adminPassword := "password123"

	regBody := map[string]string{
		"username":   adminUsername,
		"email":      adminEmail,
		"password":   adminPassword,
		"first_name": "Admin",
		"last_name":  "Reporting",
	}
	jsonRegBody, _ := json.Marshal(regBody)
	reqReg, _ := http.NewRequest("POST", "/register/", bytes.NewBuffer(jsonRegBody))
	rrReg := httptest.NewRecorder()
	router.ServeHTTP(rrReg, reqReg)

	_, err := pool.Exec(context.Background(), "UPDATE users SET role = 'admin' WHERE username = $1", adminUsername)
	if err != nil {
		t.Fatalf("Failed to update user role to admin: %v", err)
	}

	body := map[string]string{
		"username": adminUsername,
		"password": adminPassword,
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/login/", bytes.NewBuffer(jsonBody))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Failed to login as admin: %s", rr.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &resp)
	token, ok := resp["token"].(string)
	if !ok || token == "" {
		t.Fatal("No token returned in admin login")
	}
	return token
}
