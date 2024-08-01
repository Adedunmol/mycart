package services_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/models"
	"github.com/Adedunmol/mycart/internal/services"
)

type APIResponse struct {
	Message string `json:"message"`
	Data    struct {
		Token      string        `json:"token"`
		Expiration time.Duration `json:"expiration"`
	}
	Status string `json:"status"`
}

func TestMain(m *testing.M) {

	server := httptest.NewServer(http.HandlerFunc(services.CreateUserHandler))

	body := map[string]string{
		"first_name": "test",
		"last_name":  "test",
		"email":      "test@test.com",
		"username":   "testusername",
		"password":   "123456789",
	}

	postBody, _ := json.Marshal(body)

	_, err := http.Post(server.URL, "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	// drop table(s) here
	database.DB.Migrator().DropTable(&models.User{}, &models.Role{}, &models.Product{})

	os.Exit(code)
}

func TestCreateProductHandlerReturns401(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(services.CreateProductHandler))

	body := map[string]string{}

	postBody, _ := json.Marshal(body)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, but got %d", resp.StatusCode)
	}
}

func TestCreateProductHandlerReturns400(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(services.LoginUserHandler))

	body := map[string]string{
		"email":    "test@test.com",
		"password": "123456789",
	}

	postBody, _ := json.Marshal(body)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err)
	}

	var response APIResponse
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(respBody, &response)

	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()

	productBody := map[string]interface{}{
		"name":     "test",
		"details":  "some random product",
		"price":    10,
		"category": "clothing",
	}

	postProductBody, err := json.Marshal(productBody)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(postProductBody))

	if err != nil {
		t.Error(err)
	}

	services.CreateProductHandler(rr, req)

	if rr.Result().StatusCode != http.StatusUnauthorized {
		t.Errorf("expected a 401 but got %d", rr.Result().StatusCode)
	}
}

func TestCreateProductHandlerReturns200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(services.LoginUserHandler))

	body := map[string]string{
		"email":    "test@test.com",
		"password": "123456789",
	}

	postBody, _ := json.Marshal(body)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err)
	}

	var response APIResponse
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(respBody, &response)

	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()

	productBody := map[string]interface{}{
		"name":     "test",
		"details":  "some random product",
		"price":    10,
		"category": "clothing",
	}

	postProductBody, err := json.Marshal(productBody)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(postProductBody))
	rr.Header().Add("Authorization", response.Data.Token)

	if err != nil {
		t.Error(err)
	}

	services.CreateProductHandler(rr, req)

	if rr.Result().StatusCode != http.StatusCreated {
		t.Errorf("expected a 201 but got %d", rr.Result().StatusCode)
	}
}

func TestGetProductHandlerReturns400(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(services.LoginUserHandler))

	body := map[string]string{
		"email":    "test@test.com",
		"password": "123456789",
	}

	postBody, _ := json.Marshal(body)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err)
	}

	var response APIResponse
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(respBody, &response)

	if err != nil {
		t.Error(err)
	}

	server = httptest.NewServer(http.HandlerFunc(services.GetProductHandler))

	resp, err = http.Get(server.URL)

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a 400 but got %d", resp.StatusCode)
	}
}

func TestGetProductHandlerReturns404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(services.LoginUserHandler))

	body := map[string]string{
		"email":    "test@test.com",
		"password": "123456789",
	}

	postBody, _ := json.Marshal(body)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err)
	}

	var response APIResponse
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(respBody, &response)

	if err != nil {
		t.Error(err)
	}

	server = httptest.NewServer(http.HandlerFunc(services.GetProductHandler))

	resp, err = http.Get(server.URL + "/1000")

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected a 404 but got %d", resp.StatusCode)
	}
}

func TestGetProductHandlerReturns200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(services.LoginUserHandler))

	body := map[string]string{
		"email":    "test@test.com",
		"password": "123456789",
	}

	postBody, _ := json.Marshal(body)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err)
	}

	var response APIResponse
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(respBody, &response)

	if err != nil {
		t.Error(err)
	}

	server = httptest.NewServer(http.HandlerFunc(services.GetProductHandler))

	resp, err = http.Get(server.URL + "/1")

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected a 200 but got %d", resp.StatusCode)
	}
}

func TestGetAllProductsHandlerReturns400(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(services.LoginUserHandler))

	body := map[string]string{
		"email":    "test@test.com",
		"password": "123456789",
	}

	postBody, _ := json.Marshal(body)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err)
	}

	var response APIResponse
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(respBody, &response)

	if err != nil {
		t.Error(err)
	}

	server = httptest.NewServer(http.HandlerFunc(services.GetAllProductsHandler))

	resp, err = http.Get(server.URL)

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a 400 but got %d", resp.StatusCode)
	}
}

func TestGetAllProductsHandlerReturns200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(services.LoginUserHandler))

	body := map[string]string{
		"email":    "test@test.com",
		"password": "123456789",
	}

	postBody, _ := json.Marshal(body)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err)
	}

	var response APIResponse
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(respBody, &response)

	if err != nil {
		t.Error(err)
	}

	server = httptest.NewServer(http.HandlerFunc(services.GetAllProductsHandler))

	resp, err = http.Get(server.URL)

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected a 200 but got %d", resp.StatusCode)
	}
}

func TestDeleteProductHandlerReturns400(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(services.LoginUserHandler))

	body := map[string]string{
		"email":    "test@test.com",
		"password": "123456789",
	}

	postBody, _ := json.Marshal(body)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err)
	}

	var response APIResponse
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(respBody, &response)

	if err != nil {
		t.Error(err)
	}

	server = httptest.NewServer(http.HandlerFunc(services.DeleteProductHandler))

	resp, err = http.Get(server.URL)

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a 400 but got %d", resp.StatusCode)
	}
}

func TestDeleteProductHandlerReturns404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(services.LoginUserHandler))

	body := map[string]string{
		"email":    "test@test.com",
		"password": "123456789",
	}

	postBody, _ := json.Marshal(body)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err)
	}

	var response APIResponse
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(respBody, &response)

	if err != nil {
		t.Error(err)
	}

	server = httptest.NewServer(http.HandlerFunc(services.DeleteProductHandler))

	resp, err = http.Get(server.URL + "/1000")

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected a 404 but got %d", resp.StatusCode)
	}
}

func TestDeleteProductHandlerReturns403(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(services.LoginUserHandler))

	body := map[string]string{
		"email":    "test@test.com",
		"password": "123456789",
	}

	postBody, _ := json.Marshal(body)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err)
	}

	var response APIResponse
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(respBody, &response)

	if err != nil {
		t.Error(err)
	}

	server = httptest.NewServer(http.HandlerFunc(services.DeleteProductHandler))

	resp, err = http.Get(server.URL + "/1")

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected a 403 but got %d", resp.StatusCode)
	}
}

func TestDeleteProductHandlerReturns200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(services.LoginUserHandler))

	body := map[string]string{
		"email":    "test@test.com",
		"password": "123456789",
	}

	postBody, _ := json.Marshal(body)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err)
	}

	var response APIResponse
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(respBody, &response)

	if err != nil {
		t.Error(err)
	}

	server = httptest.NewServer(http.HandlerFunc(services.DeleteProductHandler))

	resp, err = http.Get(server.URL + "/1")

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected a 200 but got %d", resp.StatusCode)
	}
}

func TestUpdateProductHandlerReturns400(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(services.LoginUserHandler))

	body := map[string]string{
		"email":    "test@test.com",
		"password": "123456789",
	}

	postBody, _ := json.Marshal(body)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err)
	}

	var response APIResponse
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(respBody, &response)

	if err != nil {
		t.Error(err)
	}

	server = httptest.NewServer(http.HandlerFunc(services.UpdateProductHandler))

	resp, err = http.Get(server.URL)

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a 400 but got %d", resp.StatusCode)
	}
}

func TestUpdateProductHandlerReturns404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(services.LoginUserHandler))

	body := map[string]string{
		"email":    "test@test.com",
		"password": "123456789",
	}

	postBody, _ := json.Marshal(body)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err)
	}

	var response APIResponse
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(respBody, &response)

	if err != nil {
		t.Error(err)
	}

	server = httptest.NewServer(http.HandlerFunc(services.UpdateProductHandler))

	resp, err = http.Get(server.URL + "/1000")

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected a 404 but got %d", resp.StatusCode)
	}
}

func TestUpdateProductHandlerReturns403(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(services.LoginUserHandler))

	body := map[string]string{
		"email":    "test@test.com",
		"password": "123456789",
	}

	postBody, _ := json.Marshal(body)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err)
	}

	var response APIResponse
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(respBody, &response)

	if err != nil {
		t.Error(err)
	}

	server = httptest.NewServer(http.HandlerFunc(services.UpdateProductHandler))

	resp, err = http.Get(server.URL + "/1")

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected a 403 but got %d", resp.StatusCode)
	}
}

func TestUpdateProductHandlerReturns200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(services.LoginUserHandler))

	body := map[string]string{
		"email":    "test@test.com",
		"password": "123456789",
	}

	postBody, _ := json.Marshal(body)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err)
	}

	var response APIResponse
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(respBody, &response)

	if err != nil {
		t.Error(err)
	}

	server = httptest.NewServer(http.HandlerFunc(services.UpdateProductHandler))

	resp, err = http.Get(server.URL + "/1")

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected a 200 but got %d", resp.StatusCode)
	}
}
