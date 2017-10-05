package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeTest(t *testing.T, dockerfile string, expectedCode int, expectedBody string) {
	reqBody := bytes.NewReader([]byte(dockerfile))
	req, err := http.NewRequest("POST", "/", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(LintHandler)
	handler.ServeHTTP(rec, req)

	resCode := rec.Code
	if expectedCode != resCode {
		t.Errorf("expected %v, got %v", expectedCode, resCode)
	}

	resBody := rec.Body.String()
	if expectedBody != resBody {
		t.Errorf("expected %v, got %v", expectedBody, resBody)
	}
}

func TestContinuationWarning(t *testing.T) {
	makeTest(t, "RUN \\\n\nexit", 400, `{"error":true,"message":"Empty continuation line found in:\n    RUN exit. Empty continuation lines will become errors in a future release."}`)
}

func TestEmptyFile(t *testing.T) {
	makeTest(t, "", 400, `{"error":true,"message":"Dockerfile may not be empty"}`)
}

func TestParsingError(t *testing.T) {
	makeTest(t, "FROM", 400, `{"error":true,"message":"Dockerfile parse error line 1: FROM requires either one or three arguments"}`)
}

func TestSucess(t *testing.T) {
	makeTest(t, "FROM golang", 200, `{"error":false}`)
}
