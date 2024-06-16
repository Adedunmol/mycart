package services_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Adedunmol/mycart/internal/services"
)

// func TestMain(m *testing.M) {

// 	code := m.Run()

// 	// drop table(s) here
// 	database.DB.Migrator().DropTable(&models.User{}, &models.Role{})

// 	os.Exit(code)
// }

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

	_, err := http.Post(server.URL, "application/json", bytes.NewBuffer(postBody))

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

	resp, _ := http.Post(server.URL, "application/json", bytes.NewBuffer(postBody))

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
		fmt.Println(resp.Body)
		t.Errorf("expected 200, but got %d", resp.StatusCode)
	}

	if len(resp.Cookies()) == 0 {
		t.Error("expected a cookie to be set")
	}
}

func TestRefreshTokenHandlerReturns403(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(services.RefreshTokenHandler))

	resp, err := http.Get(server.URL)

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403, but got %d", resp.StatusCode)
	}
}

func TestRefreshTokenHandlerReturns200(t *testing.T) {
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

	var token string

	for _, c := range rr.Result().Cookies() {
		if c.Name == "token" {
			token = c.Value
		}
	}

	cookie := &http.Cookie{Name: "token", Value: token}

	req, _ = http.NewRequest(http.MethodGet, "", bytes.NewBuffer([]byte{}))
	req.AddCookie(cookie)

	services.RefreshTokenHandler(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		fmt.Println(rr.Body)
		t.Errorf("expected a 200 for login but got %d", rr.Result().StatusCode)
	}
}

func TestLogoutHandlerWithoutCookieReturns204(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(services.LogoutHandler))

	resp, err := http.Get(server.URL)

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, but got %d", resp.StatusCode)
	}
}

func TestLogoutHandlerWithCookieReturns204(t *testing.T) {
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

	var token string

	for _, c := range rr.Result().Cookies() {
		if c.Name == "token" {
			token = c.Value
		}
	}

	cookie := &http.Cookie{Name: "token", Value: token}

	req, _ = http.NewRequest(http.MethodGet, "", bytes.NewBuffer([]byte{}))
	req.AddCookie(cookie)

	services.LogoutHandler(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("expected a 204, but got %d", rr.Result().StatusCode)
	}
}
