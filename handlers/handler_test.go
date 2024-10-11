package handlers

import (
	"crproductos/db"
	"crproductos/models"
	"database/sql"
	"encoding/json"
	//"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func executeRequest(req *http.Request, s *Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

//TODO: Create test for the rest of handlers

func TestGetAllProducts(t *testing.T) {
	s := NewServer(db.ConnectToPostgres())
	s.MountHandlers()
	req := httptest.NewRequest("GET", "/products/", nil)
	response := executeRequest(req, s)
	checkResponseCode(t, http.StatusOK, response.Code)
	var result []models.ProductResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Errorf("Test failed, reason: %v", err)
	}
	var expectedProducts = []models.Product{
		{
			Id:       2,
			Name:     sql.NullString{String: "pepsi", Valid: true},
			Quantity: sql.NullFloat64{Float64: 3.000000, Valid: true},
			Unit:     sql.NullString{String: "litros", Valid: true},
			Stores:   &models.Stores{},
		},
		{
			Id:       1,
			Name:     sql.NullString{String: "coca", Valid: true},
			Quantity: sql.NullFloat64{Float64: 2.500000, Valid: true},
			Unit:     sql.NullString{String: "litros", Valid: true},
			Stores:   &models.Stores{},
		},
		{
			Id:       3,
			Name:     sql.NullString{String: "te verde", Valid: true},
			Quantity: sql.NullFloat64{Float64: 2.500000, Valid: true},
			Unit:     sql.NullString{String: "litros", Valid: true},
			Stores: &models.Stores{
				"maziplai": 3000,
				"pali":     6000,
				"walmart":  0,
			},
		},
	}
	for i, expectedValue := range expectedProducts {
		if expectedValue.ToJSON().Id != result[i].Id {
			t.Errorf("Test failed, reason: Id does not match\n Expected: %+v\n Actual: %+v\n", expectedValue.ToJSON().Id, result[i].Id)
		}
		if *expectedValue.ToJSON().Name != *result[i].Name {
			t.Errorf("Test failed, reason: Name does not match\n Expected: %+v\n Actual: %+v\n", *expectedValue.ToJSON().Name, *result[i].Name)
		}
		if *expectedValue.ToJSON().Quantity != *result[i].Quantity {
			t.Errorf("Test failed, reason: Quantity does not match\n Expected: %+v\n Actual: %+v\n", *expectedValue.ToJSON().Quantity, *result[i].Quantity)
		}
		if *expectedValue.ToJSON().Unit != *result[i].Unit {
			t.Errorf("Test failed, reason: Unit does not match\n Expected: %+v\n Actual: %+v\n", *expectedValue.ToJSON().Unit, *result[i].Unit)
		}

		for key, expectedStoreValue := range *expectedValue.ToJSON().Stores {
			if actualValue, ok := (*result[i].Stores)[key]; !ok || actualValue != expectedStoreValue {
				t.Errorf("Test failed, reason: Store value does not match\n Expected: %+v\n Actual: %+v\n", expectedStoreValue, actualValue)
			}
		}
	}
}
