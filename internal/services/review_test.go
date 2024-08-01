package services_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Adedunmol/mycart/internal/services"
)

func TestCreateReviewHandlerReturns401(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(services.CreateReviewHandler))

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

func TestCreateReviewHandlerReturns400(t *testing.T) {
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

	var response APIResponse
	respBody, err := io.ReadAll(rr.Body)

	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		t.Error(err)
	}

	reviewBody := map[string]interface{}{
		"comment": "some random comment",
		"rating":  4,
	}

	postReviewBody, err := json.Marshal(reviewBody)
	if err != nil {
		t.Error(err)
	}

	req, err = http.NewRequest(http.MethodPost, "", bytes.NewBuffer(postReviewBody))
	rr.Header().Add("Authorization", response.Data.Token)

	if err != nil {
		t.Error(err)
	}

	services.CreateReviewHandler(rr, req)

	if rr.Result().StatusCode != http.StatusCreated {
		t.Errorf("expected a 400 but got %d", rr.Result().StatusCode)
	}
}

func TestCreateReviewHandlerReturns200(t *testing.T) {
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

	reviewBody := map[string]interface{}{
		"comment": "some random comment",
		"rating":  4,
	}

	postReviewBody, err := json.Marshal(reviewBody)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(postReviewBody))
	rr.Header().Add("Authorization", response.Data.Token)

	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}

	services.CreateReviewHandler(rr, req)

	if rr.Result().StatusCode != http.StatusCreated {
		t.Errorf("expected a 201 but got %d", rr.Result().StatusCode)
	}
}

func TestGetReviewHandlerReturns400(t *testing.T) {
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

	server = httptest.NewServer(http.HandlerFunc(services.GetReviewHandler))

	resp, err = http.Get(server.URL)

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a 400 but got %d", resp.StatusCode)
	}
}

func TestGetReviewHandlerReturns404(t *testing.T) {
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

	server = httptest.NewServer(http.HandlerFunc(services.GetReviewHandler))

	resp, err = http.Get(server.URL + "/1000")

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected a 404 but got %d", resp.StatusCode)
	}
}

func TestGetReviewHandlerReturns200(t *testing.T) {
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

	server = httptest.NewServer(http.HandlerFunc(services.GetReviewHandler))

	resp, err = http.Get(server.URL + "/1")

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected a 200 but got %d", resp.StatusCode)
	}
}
