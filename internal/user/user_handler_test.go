package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

const (
	baseURL = "http://localhost:8080"
)

func TestUserCreation(t *testing.T) {
	auth := UserInput{
		Name:     "test",
		Password: "testpass",
	}
	body, _ := json.Marshal(auth)
	req, err := http.NewRequest("POST", baseURL+"/user/post/", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusConflict {
		t.Fatalf("expected 201 or 409, got %v", resp.StatusCode)
	}
}

func TestUserLogin(t *testing.T) {
	auth := UserInput{
		Name:     "test",
		Password: "testpass",
	}
	body, _ := json.Marshal(auth)
	req, err := http.NewRequest("POST", baseURL+"/user/login/", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Contet-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %v", resp.StatusCode)
	}
}
