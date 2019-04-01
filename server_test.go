package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDirectoryPathFromFilePath(t *testing.T) {
	path, fileName := getDirectoryPathFromFilePath("testDir/testSubDir/metadata.log")
	if path != "testDir/testSubDir" {
		t.Errorf(path)
	}
	if fileName != "metadata" {
		t.Errorf(fileName)
	}
}

func TestHandler(t *testing.T) {
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
