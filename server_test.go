package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Given a filepath with directories
// Then the correct content is returned
func TestGetDirectoryPathFromFilePath(t *testing.T) {
	path, fileName := getDirectoryPathFromFilePath("testDir/testSubDir/myimage.png")
	if path != "testDir/testSubDir" {
		t.Errorf(path)
	}
	if fileName != "myimage" {
		t.Errorf(fileName)
	}
}

// Given a filepath with no directories
// Then the correct content is returned
func TestGetDirectoryPathFromFilePathNoDir(t *testing.T) {
	path, fileName := getDirectoryPathFromFilePath("myimage.png")
	if path != "" {
		t.Errorf(path)
	}
	if fileName != "myimage" {
		t.Errorf(fileName)
	}
}

// When a GET request is performed
// Then the correct response is received
func TestHandlerGet(t *testing.T) {
	request, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	handler := http.HandlerFunc(handleRequest)
	handler.ServeHTTP(recorder, request)

	// Check the status code is what we expect.
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `This server sends files data to AWS S3. Note: Only multipart/form-data is accepted.`
	actual := recorder.Body.String()
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

// When an unsupported request is performed
// Then the correct response is received
func TestHandlerUnsupportedMethods(t *testing.T) {
	request, err := http.NewRequest("PUT", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	handler := http.HandlerFunc(handleRequest)
	handler.ServeHTTP(recorder, request)

	// Check the status code is what we expect.
	if status := recorder.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}

	// Check the response body is what we expect.
	expected := `Only GET and POST methods are supported.`
	actual := recorder.Body.String()
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}
