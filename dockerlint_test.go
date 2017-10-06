package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func expect(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

type testParam struct {
	method       string
	body         string
	expectedCode int
	expectedBody string
}

func makeTest(t *testing.T, param *testParam) *httptest.ResponseRecorder {
	reqBody := bytes.NewReader([]byte(param.body))
	req, err := http.NewRequest(param.method, "/", reqBody)
	expect(t, nil, err)

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(LintHandler)
	handler.ServeHTTP(rec, req)

	expect(t, param.expectedCode, rec.Code)
	expect(t, param.expectedBody, rec.Body.String())

	return rec
}

func TestNotAllowed(t *testing.T) {
	rec := makeTest(t, &testParam{
		method:       "PUT",
		body:         "",
		expectedCode: 405,
		expectedBody: "",
	})
	expect(t, "POST", rec.HeaderMap.Get("Allow"))
}

func TestContinuationWarning(t *testing.T) {
	makeTest(t, &testParam{
		method:       "POST",
		body:         "RUN \\\n\nexit",
		expectedCode: 400,
		expectedBody: `{"error":true,"message":"Empty continuation line found in:\n    RUN exit. Empty continuation lines will become errors in a future release."}`,
	})
}

func TestEmptyFile(t *testing.T) {
	makeTest(t, &testParam{
		method:       "POST",
		body:         "",
		expectedCode: 400,
		expectedBody: `{"error":true,"message":"Dockerfile may not be empty"}`,
	})
}

func TestParsingError(t *testing.T) {
	makeTest(t, &testParam{
		method:       "POST",
		body:         "FROM",
		expectedCode: 400,
		expectedBody: `{"error":true,"message":"Dockerfile parse error line 1: FROM requires either one or three arguments"}`,
	})
}

func TestSuccess(t *testing.T) {
	makeTest(t, &testParam{
		method:       "POST",
		body:         "FROM golang",
		expectedCode: 200,
		expectedBody: `{"error":false}`,
	})
}
