package services_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	body := map[string]string{
		"email":    "test@test.com",
		"password": "123456789",
	}

	postBody, _ := json.Marshal(body)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err)
	}

	services.LoginUserHandler(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("expected a 200 for login but got %d", rr.Result().StatusCode)
	}

	data, _ := io.ReadAll(rr.Result().Body)

	jsonResponse := APIResponse{}

	err = json.Unmarshal(data, &jsonResponse)

	if err != nil {
		t.Error(err)
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

	fmt.Println("token: ", response.Data.Token)

	rr := httptest.NewRecorder()

	productBody := map[string]interface{}{
		"name":     "test",
		"details":  "some random product",
		"price":    10,
		"category": "clothing",
	}

	postProductBody, err := json.Marshal(productBody)

	req, err := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(postProductBody))
	rr.Header().Add("Authorization", response.Data.Token)

	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}

	services.CreateProductHandler(rr, req)

	if rr.Result().StatusCode != http.StatusCreated {
		t.Errorf("expected a 201 but got %d", rr.Result().StatusCode)
	}
}
