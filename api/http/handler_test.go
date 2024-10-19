package http

import (
	"crproductos/internal/models"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockProductService struct{}

func (s mockProductService) GetAllProducts() ([]models.ProductResponse, error) {
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
	var result []models.ProductResponse
	for _, value := range expectedProducts {
		result = append(result, value.ToJSON())
	}
	return result, nil
}

func (s mockProductService) GetProductById(id string) (models.Product, error) {
	return models.Product{
		Id:       2,
		Name:     sql.NullString{String: "pepsi", Valid: true},
		Quantity: sql.NullFloat64{Float64: 3.000000, Valid: true},
		Unit:     sql.NullString{String: "litros", Valid: true},
		Stores:   &models.Stores{},
	}, nil
}

func (s mockProductService) CreateProduct(product models.ProductResponse) (models.ProductResponse, error) {
	var createdProduct = models.Product{
		Id:       4,
		Name:     sql.NullString{String: "te verde", Valid: true},
		Quantity: sql.NullFloat64{Float64: 2.500000, Valid: true},
		Unit:     sql.NullString{String: "litros", Valid: true},
		Stores: &models.Stores{
			"maziplai": 3000,
			"pali":     6000,
			"walmart":  0,
		},
	}
	return createdProduct.ToJSON(), nil
}
func (s mockProductService) DeleteProduct(id string) error {
	return nil
}
func (s mockProductService) UpdateProduct(id string, product models.ProductResponse) (models.ProductResponse, error) {
	return models.ProductResponse{}, nil
}
func (s mockProductService) PatchProduct(id string, product models.ProductResponse) (models.Product, error) {
	return models.Product{}, nil
}
func (s mockProductService) PatchStore(id string, jsonStore []byte) (models.Product, error) {
	return models.Product{}, nil
}

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
func TestGetAllProducts(t *testing.T) {
	s := NewServer()
	productService := &mockProductService{}
	productHandler := NewProductHandler(productService)
	s.MountHandlers(productHandler)

	req := httptest.NewRequest("GET", "/products/", nil)

	response := executeRequest(req, s)
	checkResponseCode(t, http.StatusOK, response.Code)
	fmt.Printf("Body: %v", response.Body.String())
}

//
// //TODO: Create test for the rest of handlers
//
// func TestGetAllProducts(t *testing.T) {
// 	s := NewServer()
// 	s.MountHandlers()
// 	req := httptest.NewRequest("GET", "/products/", nil)
// 	response := executeRequest(req, s)
// 	checkResponseCode(t, http.StatusOK, response.Code)
// 	var result []models.ProductResponse
// 	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
// 		t.Errorf("Test failed, reason: %v", err)
// 	}
// 	var expectedProducts = []models.Product{
// 		{
// 			Id:       2,
// 			Name:     sql.NullString{String: "pepsi", Valid: true},
// 			Quantity: sql.NullFloat64{Float64: 3.000000, Valid: true},
// 			Unit:     sql.NullString{String: "litros", Valid: true},
// 			Stores:   &models.Stores{},
// 		},
// 		{
// 			Id:       1,
// 			Name:     sql.NullString{String: "coca", Valid: true},
// 			Quantity: sql.NullFloat64{Float64: 2.500000, Valid: true},
// 			Unit:     sql.NullString{String: "litros", Valid: true},
// 			Stores:   &models.Stores{},
// 		},
// 		{
// 			Id:       3,
// 			Name:     sql.NullString{String: "te verde", Valid: true},
// 			Quantity: sql.NullFloat64{Float64: 2.500000, Valid: true},
// 			Unit:     sql.NullString{String: "litros", Valid: true},
// 			Stores: &models.Stores{
// 				"maziplai": 3000,
// 				"pali":     6000,
// 				"walmart":  0,
// 			},
// 		},
// 	}
// 	for i, expectedValue := range expectedProducts {
// 		if expectedValue.ToJSON().Id != result[i].Id {
// 			t.Errorf("Test failed, reason: Id does not match\n Expected: %+v\n Actual: %+v\n", expectedValue.ToJSON().Id, result[i].Id)
// 		}
// 		if *expectedValue.ToJSON().Name != *result[i].Name {
// 			t.Errorf("Test failed, reason: Name does not match\n Expected: %+v\n Actual: %+v\n", *expectedValue.ToJSON().Name, *result[i].Name)
// 		}
// 		if *expectedValue.ToJSON().Quantity != *result[i].Quantity {
// 			t.Errorf("Test failed, reason: Quantity does not match\n Expected: %+v\n Actual: %+v\n", *expectedValue.ToJSON().Quantity, *result[i].Quantity)
// 		}
// 		if *expectedValue.ToJSON().Unit != *result[i].Unit {
// 			t.Errorf("Test failed, reason: Unit does not match\n Expected: %+v\n Actual: %+v\n", *expectedValue.ToJSON().Unit, *result[i].Unit)
// 		}
//
// 		for key, expectedStoreValue := range *expectedValue.ToJSON().Stores {
// 			if actualValue, ok := (*result[i].Stores)[key]; !ok || actualValue != expectedStoreValue {
// 				t.Errorf("Test failed, reason: Store value does not match\n Expected: %+v\n Actual: %+v\n", expectedStoreValue, actualValue)
// 			}
// 		}
// 	}
// }
//
// func TestGetProductById(t *testing.T) {
// 	s := NewServer(db.ConnectToPostgres())
// 	s.MountHandlers()
// 	req := httptest.NewRequest("GET", "/products/2", nil)
// 	response := executeRequest(req, s)
// 	checkResponseCode(t, http.StatusOK, response.Code)
// 	var result models.ProductResponse
// 	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
// 		t.Errorf("Test failed, reason: %v", err)
// 	}
// 	var expectedValue = models.Product{
// 		Id:       2,
// 		Name:     sql.NullString{String: "pepsi", Valid: true},
// 		Quantity: sql.NullFloat64{Float64: 3.000000, Valid: true},
// 		Unit:     sql.NullString{String: "litros", Valid: true},
// 		Stores:   &models.Stores{},
// 	}
// 	if expectedValue.ToJSON().Id != result.Id {
// 		t.Errorf("Test failed, reason: Id does not match\n Expected: %+v\n Actual: %+v\n", expectedValue.ToJSON().Id, result.Id)
// 	}
// 	if *expectedValue.ToJSON().Name != *result.Name {
// 		t.Errorf("Test failed, reason: Name does not match\n Expected: %+v\n Actual: %+v\n", *expectedValue.ToJSON().Name, *result.Name)
// 	}
// 	if *expectedValue.ToJSON().Quantity != *result.Quantity {
// 		t.Errorf("Test failed, reason: Quantity does not match\n Expected: %+v\n Actual: %+v\n", *expectedValue.ToJSON().Quantity, *result.Quantity)
// 	}
// 	if *expectedValue.ToJSON().Unit != *result.Unit {
// 		t.Errorf("Test failed, reason: Unit does not match\n Expected: %+v\n Actual: %+v\n", *expectedValue.ToJSON().Unit, *result.Unit)
// 	}
//
// 	for key, expectedStoreValue := range *expectedValue.ToJSON().Stores {
// 		if actualValue, ok := (*result.Stores)[key]; !ok || actualValue != expectedStoreValue {
// 			t.Errorf("Test failed, reason: Store value does not match\n Expected: %+v\n Actual: %+v\n", expectedStoreValue, actualValue)
// 		}
// 	}
// }
// func TestGetProductByIdFail(t *testing.T) {
// 	s := NewServer(db.ConnectToPostgres())
// 	s.MountHandlers()
// 	req := httptest.NewRequest("GET", "/products/undefined", nil)
// 	response := executeRequest(req, s)
//
// 	checkResponseCode(t, http.StatusInternalServerError, response.Code)
// }
//
// var createdID int
//
// func TestCreateProduct(t *testing.T) {
// 	var requestBody = models.Product{
// 		Name:     sql.NullString{String: "te verde", Valid: true},
// 		Quantity: sql.NullFloat64{Float64: 2.500000, Valid: true},
// 		Unit:     sql.NullString{String: "litros", Valid: true},
// 		Stores: &models.Stores{
// 			"maziplai": 3000,
// 			"pali":     6000,
// 			"walmart":  0,
// 		},
// 	}
// 	s := NewServer(db.ConnectToPostgres())
// 	s.MountHandlers()
// 	b, err := json.Marshal(requestBody.ToJSON())
// 	if err != nil {
// 		t.Fatalf("Could not create request: %+v", err)
// 	}
// 	req := httptest.NewRequest("POST", "/products/", bytes.NewBuffer(b))
// 	response := executeRequest(req, s)
// 	checkResponseCode(t, http.StatusOK, response.Code)
// 	var result models.ProductResponse
// 	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
// 		t.Errorf("Test failed, reason: %v", err)
// 	}
// 	if *requestBody.ToJSON().Name != *result.Name {
// 		t.Errorf("Test failed, reason: Name does not match\n Expected: %+v\n Actual: %+v\n", *requestBody.ToJSON().Name, *result.Name)
// 	}
// 	if *requestBody.ToJSON().Quantity != *result.Quantity {
// 		t.Errorf("Test failed, reason: Quantity does not match\n Expected: %+v\n Actual: %+v\n", *requestBody.ToJSON().Quantity, *result.Quantity)
// 	}
// 	if *requestBody.ToJSON().Unit != *result.Unit {
// 		t.Errorf("Test failed, reason: Unit does not match\n Expected: %+v\n Actual: %+v\n", *requestBody.ToJSON().Unit, *result.Unit)
// 	}
//
// 	for key, expectedStoreValue := range *requestBody.ToJSON().Stores {
// 		if actualValue, ok := (*result.Stores)[key]; !ok || actualValue != expectedStoreValue {
// 			t.Errorf("Test failed, reason: Store value does not match\n Expected: %+v\n Actual: %+v\n", expectedStoreValue, actualValue)
// 		}
// 	}
// 	createdID = result.Id
// }
//
// func TestDeleteProduct(t *testing.T) {
// 	s := NewServer(db.ConnectToPostgres())
// 	s.MountHandlers()
// 	req := httptest.NewRequest("DELETE", fmt.Sprintf("/products/%v", createdID), nil)
// 	response := executeRequest(req, s)
// 	checkResponseCode(t, http.StatusOK, response.Code)
//
// }
