package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCompaniesAPI(t *testing.T) {
	pool := GetPool()
	router := SetupTestRouter(pool)

	validSiret := "44306184100047"
	invalidSiret := "12345678901234"
	wrongLengthSiret := "123"

	var createdCompanyId int64

	t.Run("Create_ValidCompany", func(t *testing.T) {
		reqGet, _ := http.NewRequest("GET", "/companies/siret?siret="+validSiret, nil)
		rrGet := httptest.NewRecorder()
		router.ServeHTTP(rrGet, reqGet)

		if rrGet.Code == http.StatusOK {
			var existing map[string]interface{}
			json.Unmarshal(rrGet.Body.Bytes(), &existing)
			createdCompanyId = int64(existing["id"].(float64))
			t.Log("Company already exists, skipping creation")
			return
		}

		body := map[string]string{
			"siret":   validSiret,
			"name":    "Google France",
			"address": "8 Rue de Londres, 75009 Paris",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/companies", bytes.NewBuffer(jsonBody))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status Created, got %v: %s", rr.Code, rr.Body.String())
			return
		}

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		createdCompanyId = int64(resp["id"].(float64))
	})

	t.Run("Create_ExistingCompany", func(t *testing.T) {
		body := map[string]string{
			"siret":   validSiret,
			"name":    "Duplicate",
			"address": "Somewhere",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/companies", bytes.NewBuffer(jsonBody))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status BadRequest for duplicate, got %v", rr.Code)
		}
	})

	t.Run("Create_InvalidSiretLength", func(t *testing.T) {
		body := map[string]string{
			"siret": wrongLengthSiret,
			"name":  "Invalid",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/companies", bytes.NewBuffer(jsonBody))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status BadRequest for wrong length, got %v", rr.Code)
		}
	})

	t.Run("Create_NonExistentSiret", func(t *testing.T) {
		body := map[string]string{
			"siret": invalidSiret,
			"name":  "Non Existent",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/companies", bytes.NewBuffer(jsonBody))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status BadRequest for non-existent SIRET, got %v", rr.Code)
		}
	})

	t.Run("Create_InvalidJSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/companies", bytes.NewBufferString("{invalid json"))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status BadRequest for invalid JSON, got %v", rr.Code)
		}
	})

	t.Run("GetCompany_ById_Success", func(t *testing.T) {
		if createdCompanyId == 0 {
			t.Skip("No company created")
		}
		req, _ := http.NewRequest("GET", fmt.Sprintf("/companies/%d", createdCompanyId), nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", rr.Code)
		}
	})

	t.Run("GetCompany_ById_NotFound", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/companies/999999", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status NotFound, got %v", rr.Code)
		}
	})

	t.Run("GetCompany_BySiret_Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/companies/siret?siret="+validSiret, nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", rr.Code)
		}
	})

	t.Run("GetCompany_BySiret_NotFound", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/companies/siret?siret=00000000000000", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status NotFound, got %v", rr.Code)
		}
	})
}
