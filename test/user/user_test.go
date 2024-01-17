package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Adedunmol/mycart/internal/app"
	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/models"
	"github.com/Adedunmol/mycart/internal/services"
)

func TestMain(m *testing.M) {

	_ = app.Initializers()

	code := m.Run()

	// drop table(s) here
	database.Database.DB.Migrator().DropTable(&models.User{}, &models.Role{})

	os.Exit(code)
}

func TestCreateUserHandlerReturns201(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(services.CreateUserHandler))

	body := map[string]string{
		"first_name": "test",
		"last_name":  "test",
		"email":      "test@test.com",
		"username":   "testusername",
		"password":   "123456789",
	}

	postBody, _ := json.Marshal(body)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, but got %d", resp.StatusCode)
	}
}

func TestCreateUserHandlerReturns409(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(services.CreateUserHandler))

	body := map[string]string{
		"first_name": "test",
		"last_name":  "test",
		"email":      "test@test.com",
		"username":   "testusername",
		"password":   "123456789",
	}

	postBody, _ := json.Marshal(body)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err)
	}

	body = map[string]string{
		"first_name": "test",
		"last_name":  "test",
		"email":      "test@test.com",
		"username":   "testusername",
		"password":   "123456789",
	}

	postBody, _ = json.Marshal(body)

	resp, err = http.Post(server.URL, "application/json", bytes.NewBuffer(postBody))

	if resp.StatusCode != http.StatusConflict {
		t.Errorf("expected 409, but got %d", resp.StatusCode)
	}
}

func TestLoginUserHandlerReturns400(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(services.LoginUserHandler))

	body := map[string]string{}

	postBody, _ := json.Marshal(body)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, but got %d", resp.StatusCode)
	}
}

func TestLoginUserHandlerReturns401(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(services.LoginUserHandler))

	body := map[string]string{
		"email":    "test@test.com",
		"password": "123456",
	}

	postBody, _ := json.Marshal(body)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, but got %d", resp.StatusCode)
	}
}

func TestLoginUserHandlerReturns200(t *testing.T) {
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

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, but got %d", resp.StatusCode)
	}
}
