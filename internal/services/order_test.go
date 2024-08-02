package services_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Adedunmol/mycart/internal/services"
)

func TestCreateOrderReturns400(t *testing.T) {
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

	// rr := httptest.NewRecorder()

	server = httptest.NewServer(http.HandlerFunc(services.CreateOrderHandler))

	resp, err = http.Get(server.URL)

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a 400 but got %d", resp.StatusCode)
	}
}

func TestCreateOrderReturns401(t *testing.T) {
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

	// rr := httptest.NewRecorder()

	server = httptest.NewServer(http.HandlerFunc(services.CreateOrderHandler))

	resp, err = http.Get(server.URL)

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected a 401 but got %d", resp.StatusCode)
	}
}

func TestCreateOrderReturns200(t *testing.T) {
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

	// rr := httptest.NewRecorder()

	server = httptest.NewServer(http.HandlerFunc(services.CreateOrderHandler))

	resp, err = http.Get(server.URL)

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected a 200 but got %d", resp.StatusCode)
	}
}
