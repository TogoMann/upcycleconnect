package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestReportingPDF(t *testing.T) {
	pool := GetPool()
	router := SetupTestRouter(pool)
	adminToken := getAdminTokenForPDF(t, router, pool)

	t.Run("ExportAuditPDF_Success", func(t *testing.T) {
		start := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
		end := time.Now().Format("2006-01-02")
		url := "/reporting/audit/items/pdf?start=" + start + "&end=" + end

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d: %s", rr.Code, rr.Body.String())
		}

		contentType := rr.Header().Get("Content-Type")
		if contentType != "application/pdf" {
			t.Errorf("expected content-type application/pdf, got %s", contentType)
		}
	})

	t.Run("ExportAuditPDF_MissingParams", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/reporting/audit/items/pdf", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rr.Code)
		}
	})

	t.Run("ExportAuditPDF_InvalidDatesChronology", func(t *testing.T) {
		start := "2026-12-31"
		end := "2026-01-01"
		url := "/reporting/audit/items/pdf?start=" + start + "&end=" + end

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rr.Code)
		}
	})
}

func getAdminTokenForPDF(t *testing.T, router http.Handler, pool *pgxpool.Pool) string {
	adminUsername := "admin_pdf_test"
	adminEmail := "admin_pdf_test@example.com"
	adminPassword := "password123"

	regBody := map[string]string{
		"username":   adminUsername,
		"email":      adminEmail,
		"password":   adminPassword,
		"first_name": "Admin",
		"last_name":  "PDF",
	}
	jsonRegBody, _ := json.Marshal(regBody)
	reqReg, _ := http.NewRequest("POST", "/register/", bytes.NewBuffer(jsonRegBody))
	rrReg := httptest.NewRecorder()
	router.ServeHTTP(rrReg, reqReg)

	_, err := pool.Exec(context.Background(), "UPDATE users SET role = 'admin' WHERE username = $1", adminUsername)
	if err != nil {
		t.Fatalf("Failed to update user to admin: %v", err)
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
	return resp["token"].(string)
}
