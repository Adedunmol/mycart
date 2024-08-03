package services_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestAddToCartReturns400(t *testing.T) {
	clearTables()
	user := createUser()
	token, _ := generateToken(user.Username, time.Duration(15))

	body := map[string]string{}

	postBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/carts/", bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", token)

	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestAddToCartReturns404(t *testing.T) {
	clearTables()
	user := createUser()
	token, _ := generateToken(user.Username, time.Duration(15))

	body := map[string]string{}

	postBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/carts/?product_id=100&quantity=1", bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", token)

	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestAddToCartReturns200(t *testing.T) {
	clearTables()
	user := createUser()
	product, _ := createProduct()
	token, _ := generateToken(user.Username, time.Duration(15))

	body := map[string]string{}

	postBody, _ := json.Marshal(body)

	productID := strconv.Itoa(int(product.ID))
	req, _ := http.NewRequest("POST", "/carts/?product_id="+productID+"&quantity=1", bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", token)

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetCartReturns200(t *testing.T) {
	clearTables()
	user := createUser()
	addItemToCart(int(user.ID), 0)
	token, _ := generateToken(user.Username, time.Duration(15))

	body := map[string]string{}

	postBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("GET", "/carts/", bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", token)

	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestRemoveFromCartReturns400(t *testing.T) {
	clearTables()
	user := createUser()
	token, _ := generateToken(user.Username, time.Duration(15))

	body := map[string]string{}

	postBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/carts/", bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", token)

	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestRemoveCartReturns200(t *testing.T) {
	clearTables()
	user := createUser()
	product, _ := createProduct()
	addItemToCart(int(user.ID), int(product.ID))
	token, _ := generateToken(user.Username, time.Duration(15))

	body := map[string]string{}

	postBody, _ := json.Marshal(body)

	productID := strconv.Itoa(int(product.ID))
	req, _ := http.NewRequest("POST", "/carts/?product_id="+productID+"&quantity=1", bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", token)

	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}
