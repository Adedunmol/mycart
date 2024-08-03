package services_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestCreateProductHandlerReturns401(t *testing.T) {
	clearTables()

	body := map[string]string{}

	postBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

func TestCreateProductHandlerReturns422(t *testing.T) {
	clearTables()
	user := createVendor()
	token, _ := generateToken(user.Username, time.Duration(15))

	body := map[string]string{}

	postBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+strings.TrimSpace(token))

	response := executeRequest(req)

	checkResponseCode(t, http.StatusUnprocessableEntity, response.Code)
}

func TestCreateProductHandlerReturns403(t *testing.T) {
	clearTables()
	user := createUser()
	token, _ := generateToken(user.Username, time.Duration(15))

	body := map[string]interface{}{
		"name":     "test",
		"details":  "some random product",
		"price":    10,
		"category": "clothing",
		"Quantity": 100,
	}
	postBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+strings.TrimSpace(token))

	response := executeRequest(req)

	checkResponseCode(t, http.StatusForbidden, response.Code)
}

func TestCreateProductHandlerReturns201(t *testing.T) {
	clearTables()
	user := createVendor()
	token, _ := generateToken(user.Username, time.Duration(15))

	body := map[string]interface{}{
		"name":     "test",
		"details":  "some random product",
		"price":    10,
		"category": "clothing",
		"quantity": 100,
	}
	postBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+strings.TrimSpace(token))

	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)
}

func TestGetProductHandlerReturns404(t *testing.T) {
	clearTables()
	user := createUser()
	token, _ := generateToken(user.Username, time.Duration(15))

	req, _ := http.NewRequest("GET", "/products/100", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+strings.TrimSpace(token))

	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestGetProductHandlerReturns200(t *testing.T) {
	clearTables()
	user := createUser()
	product, _ := createProduct()
	token, _ := generateToken(user.Username, time.Duration(15))

	productID := strconv.Itoa(int(product.ID))
	req, _ := http.NewRequest("GET", "/products/"+productID, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+strings.TrimSpace(token))

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetAllProductsHandlerReturns200(t *testing.T) {
	clearTables()
	user := createUser()
	createProduct()
	token, _ := generateToken(user.Username, time.Duration(15))

	req, _ := http.NewRequest("GET", "/products/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+strings.TrimSpace(token))

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestDeleteProductHandlerReturns404(t *testing.T) {
	clearTables()
	user := createVendor()
	token, _ := generateToken(user.Username, time.Duration(15))

	req, _ := http.NewRequest("DELETE", "/products/100", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+strings.TrimSpace(token))

	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestDeleteProductHandlerReturns403(t *testing.T) {
	clearTables()
	user := createUser()
	product, _ := createProduct()
	token, _ := generateToken(user.Username, time.Duration(15))

	productID := strconv.Itoa(int(product.ID))
	req, _ := http.NewRequest("DELETE", "/products/"+productID, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+strings.TrimSpace(token))

	response := executeRequest(req)

	checkResponseCode(t, http.StatusForbidden, response.Code)
}

func TestDeleteProductHandlerReturns200(t *testing.T) {
	clearTables()
	product, user := createProduct()
	token, _ := generateToken(user.Username, time.Duration(15))

	productID := strconv.Itoa(int(product.ID))
	req, _ := http.NewRequest("DELETE", "/products/"+productID, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+strings.TrimSpace(token))

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateProductHandlerReturns404(t *testing.T) {
	clearTables()
	_, user := createProduct()
	token, _ := generateToken(user.Username, time.Duration(15))

	req, _ := http.NewRequest("PATCH", "/products/100", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+strings.TrimSpace(token))

	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestUpdateProductHandlerReturns403(t *testing.T) {
	clearTables()
	user := createUser()
	product, _ := createProduct()
	token, _ := generateToken(user.Username, time.Duration(15))

	productID := strconv.Itoa(int(product.ID))
	req, _ := http.NewRequest("PATCH", "/products/"+productID, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+strings.TrimSpace(token))

	response := executeRequest(req)

	checkResponseCode(t, http.StatusForbidden, response.Code)
}

func TestUpdateProductHandlerReturns200(t *testing.T) {
	clearTables()
	product, user := createProduct()
	token, _ := generateToken(user.Username, time.Duration(15))

	body := map[string]interface{}{
		"name":     "test",
		"details":  "some random product",
		"price":    10,
		"category": "clothing",
		"quantity": 100,
	}
	postBody, _ := json.Marshal(body)

	productID := strconv.Itoa(int(product.ID))
	req, _ := http.NewRequest("PATCH", "/products/"+productID, bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+strings.TrimSpace(token))

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}
