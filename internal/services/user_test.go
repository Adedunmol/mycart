package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUserHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(CreateUserHandler))

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
