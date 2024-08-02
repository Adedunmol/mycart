package services_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestCreateReviewHandlerReturns401(t *testing.T) {
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

	body := map[string]string{}

	postBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/reviews/", bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&cookie)

	response := executeRequest(req)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

func TestCreateReviewHandlerReturns400(t *testing.T) {
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

	body := map[string]string{}

	postBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/reviews/", bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&cookie)

	response := executeRequest(req)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)

}

func TestCreateReviewHandlerReturns200(t *testing.T) {
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

	reviewBody := map[string]interface{}{
		"comment": "some random comment",
		"rating":  4,
	}

	postReviewBody, err := json.Marshal(reviewBody)
	if err != nil {
		t.Error(err)
	}

	req, _ := http.NewRequest("POST", "/reviews/", bytes.NewBuffer(postReviewBody))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&cookie)
	req.Header.Add("Authorization", token)

	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)
}

func TestGetReviewHandlerReturns400(t *testing.T) {
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

	req, _ := http.NewRequest("GET", "/reviews/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&cookie)
	req.Header.Add("Authorization", token)

	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)
}

func TestGetReviewHandlerReturns404(t *testing.T) {
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

	reviewBody := map[string]interface{}{
		"comment": "some random comment",
		"rating":  4,
	}

	postReviewBody, err := json.Marshal(reviewBody)
	if err != nil {
		t.Error(err)
	}

	req, _ := http.NewRequest("POST", "/reviews/100", bytes.NewBuffer(postReviewBody))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&cookie)
	req.Header.Add("Authorization", token)

	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestGetReviewHandlerReturns200(t *testing.T) {
	clearTables()
	user := createUser()
	token, _ := generateToken(user.Username, time.Duration(15))
	product := createProduct()

	cookie := http.Cookie{
		Name:  "token",
		Value: token,
		// Expires:  time.Now().Add(util.REFRESH_TOKEN_EXPIRATION),
		HttpOnly: true,
		MaxAge:   1 * 60 * 60,
	}

	reviewBody := map[string]interface{}{
		"comment": "some random comment",
		"rating":  4,
	}

	postReviewBody, err := json.Marshal(reviewBody)
	if err != nil {
		t.Error(err)
	}
	productID := strconv.Itoa(int(product.ID))

	req, _ := http.NewRequest("POST", "/reviews/"+productID, bytes.NewBuffer(postReviewBody))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&cookie)
	req.Header.Add("Authorization", token)

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}
