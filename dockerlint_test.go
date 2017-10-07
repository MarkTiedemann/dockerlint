package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMain(t *testing.T) {
	go func() {
		main()
	}()
}

func Test404(t *testing.T) {
	makeTest(t, &Param{
		path:         "/invalid-path",
		expectedCode: 404,
	})
}

func Test405(t *testing.T) {
	rec := makeTest(t, &Param{
		method:       "PUT",
		expectedCode: 405,
	})
	equal(t, "POST", rec.HeaderMap.Get("Allow"))
}

func Test400ContinuationWarning(t *testing.T) {
	makeTest(t, &Param{
		body:         "RUN \\\n\nexit",
		expectedCode: 400,
		expectedBody: `{"error":true,"message":"Empty continuation line found in:\n    RUN exit. Empty continuation lines will become errors in a future release."}`,
	})
}

func Test400EmptyFile(t *testing.T) {
	makeTest(t, &Param{
		expectedCode: 400,
		expectedBody: `{"error":true,"message":"Dockerfile may not be empty"}`,
	})
}

func Test400ParsingError(t *testing.T) {
	makeTest(t, &Param{
		body:         "FROM",
		expectedCode: 400,
		expectedBody: `{"error":true,"message":"Dockerfile parse error line 1: FROM requires either one or three arguments"}`,
	})
}

func Test200(t *testing.T) {
	makeTest(t, &Param{
		body:         "FROM golang",
		expectedBody: `{"error":false}`,
	})
}

// UTILS

func equal(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

type Param struct {
	path         string
	method       string
	body         string
	expectedCode int
	expectedBody string
}

func newParam(param *Param) *Param {
	if param.path == "" {
		param.path = "/"
	}
	if param.method == "" {
		param.method = "POST"
	}
	if param.expectedCode == 0 {
		param.expectedCode = 200
	}
	return param
}

func makeTest(t *testing.T, param *Param) *httptest.ResponseRecorder {
	param = newParam(param)

	reqBody := bytes.NewReader([]byte(param.body))
	req, err := http.NewRequest(param.method, param.path, reqBody)
	equal(t, nil, err)

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(lintHandler)
	handler.ServeHTTP(rec, req)

	equal(t, param.expectedCode, rec.Code)
	equal(t, param.expectedBody, rec.Body.String())

	return rec
}
