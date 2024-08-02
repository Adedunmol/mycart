package services_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestCreateUserHandlerReturns201(t *testing.T) {
	clearTables()

	body := map[string]string{
		"first_name": "test",
		"last_name":  "test",
		"email":      "test@test.com",
		"username":   "testusername",
		"password":   "123456789",
	}

	postBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer([]byte(postBody)))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

}

func TestCreateUserHandlerReturns409(t *testing.T) {
	clearTables()
	createUser()

	body := map[string]string{
		"first_name": "test",
		"last_name":  "test",
		"email":      "test@test.com",
		"username":   "testusername",
		"password":   "123456789",
	}

	postBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer([]byte(postBody)))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusConflict, response.Code)
}

func TestLoginUserHandlerReturns400(t *testing.T) {
	clearTables()

	body := map[string]string{}

	postBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer([]byte(postBody)))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

}

func TestLoginUserHandlerReturns401(t *testing.T) {
	clearTables()
	createUser()

	body := map[string]string{
		"email":    "test@test.com",
		"password": "1234567",
	}

	postBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer([]byte(postBody)))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

func TestLoginUserHandlerReturns200(t *testing.T) {
	clearTables()
	createUser()

	body := map[string]string{
		"email":    "test@test.com",
		"password": "123456789",
	}

	postBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer([]byte(postBody)))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if len(response.Result().Cookies()) == 0 {
		t.Error("expected a cookie to be set")
	}
}

func TestRefreshTokenHandlerReturns403(t *testing.T) {
	clearTables()
	createUser()

	req, _ := http.NewRequest("GET", "/users/refresh", nil)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusForbidden, response.Code)
}

func TestRefreshTokenHandlerReturns200(t *testing.T) {
	clearTables()
	user := createUser()
	token, _ := generateToken(user.Username, time.Duration(15))

	cookie := http.Cookie{
		Name:  "token",
		Value: token,
		// Expires:  time.Now().Add(util.REFRESH_TOKEN_EXPIRATION),
		HttpOnly: true,
		MaxAge:   1 * 60 * 60,
	}

	req, _ := http.NewRequest("GET", "/users/refresh", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&cookie)

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestLogoutHandlerWithoutCookieReturns204(t *testing.T) {
	clearTables()
	createUser()

	req, _ := http.NewRequest("GET", "/users/logout", nil)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusNoContent, response.Code)
}

func TestLogoutHandlerWithCookieReturns204(t *testing.T) {
	clearTables()
	user := createUser()
	token, _ := generateToken(user.Username, time.Duration(15))

	cookie := http.Cookie{
		Name:  "token",
		Value: token,
		// Expires:  time.Now().Add(util.REFRESH_TOKEN_EXPIRATION),
		HttpOnly: true,
		MaxAge:   1 * 60 * 60,
	}

	req, _ := http.NewRequest("GET", "/users/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&cookie)

	response := executeRequest(req)

	checkResponseCode(t, http.StatusNoContent, response.Code)
}
