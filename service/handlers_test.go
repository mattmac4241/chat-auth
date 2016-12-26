package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var (
	formatter = render.New(render.Options{
		IndentJSON: true,
	})
)

func TestRegisterUserHandlerInvalidJSON(t *testing.T) {
	database := &testDatabase{}

	client := &http.Client{}
	server := httptest.NewServer(http.HandlerFunc(registerUserHandler(formatter, database)))
	defer server.Close()

	body := []byte("this is not valid json")
	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))

	if err != nil {
		t.Errorf("Error in creating POST request for registerUserHandler: %v", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to registerUserHandler: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Error("Sending invalid JSON should result in a bad request from server.")
	}
}

func TestRegisterUserHandlerNotUser(t *testing.T) {
	database := &testDatabase{}
	client := &http.Client{}
	server := httptest.NewServer(http.HandlerFunc(registerUserHandler(formatter, database)))
	defer server.Close()

	body := []byte("{\"test\":\"Not comment.\"}")
	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))

	if err != nil {
		t.Errorf("Error in creating second POST request for invalid data on create user: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	res, _ := client.Do(req)
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Error("Sending valid JSON but with incorrect or missing fields should result in a bad request and didn't.")
	}
}

func TestRegisterUserHandlerSuccess(t *testing.T) {
	database := &testDatabase{}
	client := &http.Client{}
	server := httptest.NewServer(http.HandlerFunc(registerUserHandler(formatter, database)))
	defer server.Close()

	body := []byte("{\"username\": \"testname\",\n \"password\": \"password\",\n \"email\": \"test@mail.com\"}")
	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in creating second POST request for invalid data on create user: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusBadRequest {
		t.Error("Sending valid JSON but with incorrect or missing fields should result in a bad request and didn't.")
	}
	if len(database.users) <= 0 || len(database.tokens) <= 0 {
		t.Error("User and token not added to database")
	}
}

func TestLoginUserHandlerInvalidJSON(t *testing.T) {
	database := &testDatabase{}

	client := &http.Client{}
	server := httptest.NewServer(http.HandlerFunc(loginUserHandler(formatter, database)))
	defer server.Close()

	body := []byte("this is not valid json")
	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))

	if err != nil {
		t.Errorf("Error in creating POST request for registerUserHandler: %v", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to registerUserHandler: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Error("Sending invalid JSON should result in a bad request from server.")
	}
}

func TestLoginUserHandlerNotUser(t *testing.T) {
	database := &testDatabase{}
	client := &http.Client{}
	server := httptest.NewServer(http.HandlerFunc(loginUserHandler(formatter, database)))
	defer server.Close()

	body := []byte("{\"test\":\"Not comment.\"}")
	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))

	if err != nil {
		t.Errorf("Error in creating second POST request for invalid data on create user: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	res, _ := client.Do(req)
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Error("Sending valid JSON but with incorrect or missing fields should result in a bad request and didn't.")
	}
}

func TestLoginUserHandlerSuccess(t *testing.T) {
	user := User{Username: "testname", Password: "password", Email: "test@mail.com", ID: 1}
	database := &testDatabase{}
	user.Save(database)
	client := &http.Client{}
	server := httptest.NewServer(http.HandlerFunc(loginUserHandler(formatter, database)))
	defer server.Close()

	body := []byte("{\"username\":\"testname\",\n\"password\":\"password\"}")
	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))

	if err != nil {
		t.Errorf("Error in creating second POST request for invalid data on create user: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	res, _ := client.Do(req)
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Error("Sending valid user shouldn't result in a bad request and didn't.")
	}
}

func TestGetTokenValidationInvalidToken(t *testing.T) {
	var (
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)
	database := &testDatabase{}

	server := MakeTestServer(database)
	token := &Token{Key: "test", ExpiresAt: time.Now().Unix()}
	database.addToken(token)

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "/auth/token/test2", nil)
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("Expected %v; received %v", http.StatusNotFound, recorder.Code)
	}
}

func TestGetTokenValidationValidToken(t *testing.T) {
	var (
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)
	database := &testDatabase{}

	server := MakeTestServer(database)
	token := &Token{Key: "test", ExpiresAt: time.Now().Add(time.Hour * 24 * 7 * time.Duration(8)).Unix()}
	token.Save(database)

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "/auth/token/"+token.Key, nil)
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected %v; received %v", http.StatusOK, recorder.Code)
	}

	var tokenResponse Token
	err := json.Unmarshal(recorder.Body.Bytes(), &tokenResponse)
	if err != nil {
		t.Errorf("Error unmarshaling token: %s", err)
	}
	if tokenResponse.Key != "test" {
		t.Errorf("Expected token key to be test; received %s", token.Key)
	}
}

func MakeTestServer(database *testDatabase) *negroni.Negroni {
	server := negroni.New()
	mx := mux.NewRouter()
	initRoutes(mx, formatter, database)
	server.UseHandler(mx)
	return server
}
