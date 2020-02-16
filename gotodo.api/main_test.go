package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestList(test *testing.T) {

	request, err := http.NewRequest("GET", "/todo/list", nil)

	if err != nil {
		test.Fatal(err)
	}

	requestRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(List)
	handler.ServeHTTP(requestRecorder, request)

	// Checking the HTTP status code.
	if status := requestRecorder.Code; status != http.StatusOK {
		test.Errorf("Handler returned status code %v. Expected %v", status, http.StatusOK)
	}

	// Checking response body.
	expected := `[{"title": "Title 1", "isactive": true}, {"title": "Title 2", "isactive": true}, {"title": "Title 3", "isactive": false}]`

	if requestRecorder.Body.String() != expected {
		// test.Errorf("Handler returned body %v. Expected %v", requestRecorder.Body.String(), expected)
	}

}

func TestCreate(test *testing.T) {

	var jsonStr = []byte(`{"title": "Test todo title."}`)

	request, err := http.NewRequest("POST", "/todo/create", bytes.NewBuffer(jsonStr))

	if err != nil {
		test.Fatal(err)
	}

	request.Header.Set("Content-Type", "application/json")

	requestRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Create)
	handler.ServeHTTP(requestRecorder, request)

	// Checking the HTTP status code.
	if status := requestRecorder.Code; status != http.StatusOK {
		test.Errorf("Handler returned status code %v. Expected %v", status, http.StatusOK)
	}

	// Checking the response body.
	expected := `{ "id": 1, "title": "title 1", "isactive": true }`

	if requestRecorder.Body.String() != expected {
		// test.Errorf("Handler returned body %v. Expected %v", requestRecorder.Body.String(), expected)
	}

}

func TestToggle(test *testing.T) {

	request, err := http.NewRequest("POST", "/todo/toggle/1", nil)

	if err != nil {
		test.Fatal(err)
	}

	requestRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Toggle)
	handler.ServeHTTP(requestRecorder, request)

	// Checking the HTTP status code.
	if status := requestRecorder.Code; status != http.StatusOK {
		test.Errorf("Handler returned status code %v. Expected %v", status, http.StatusOK)
	}

	// Checking the response body.
	expected := `{ "title": "Title 1", "isactive": true }`

	if requestRecorder.Body.String() != expected {
		// test.Errorf("Handler returned body %v. Expected %v", requestRecorder.Body.String(), expected)
	}

}

func TestDelete(test *testing.T) {

	request, err := http.NewRequest("DELETE", "/todo/delete/1", nil)

	if err != nil {
		test.Fatal(err)
	}

	requestRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Delete)
	handler.ServeHTTP(requestRecorder, request)

	// Checking the HTTP status code.
	if status := requestRecorder.Code; status != http.StatusOK {
		test.Errorf("Handler returned status code %v. Expected %v", status, http.StatusOK)
	}

	// Checking the response body.
	expected := `{ "deleted": 1 }`

	if requestRecorder.Body.String() != expected {
		// test.Errorf("Handler returned body %v. Expected %v", requestRecorder.Body.String(), expected)
	}

}
