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
	makeTest(t, &param{
		path:         "/invalid-path",
		expectedCode: 404,
	})
}

func Test405(t *testing.T) {
	rec := makeTest(t, &param{
		method:       "PUT",
		expectedCode: 405,
	})
	equal(t, "POST", rec.HeaderMap.Get("Allow"))
}

func Test400ParseError(t *testing.T) {
	makeTest(t, &param{
		body:         "FROM busybox\n\nENV PATH",
		expectedCode: 400,
		expectedBody: `{"error":true,"message":"ENV must have two arguments"}`,
	})
}

func Test400ContinuationWarning(t *testing.T) {
	makeTest(t, &param{
		body:         "RUN \\\n\nexit",
		expectedCode: 400,
		expectedBody: `{"error":true,"message":"Empty continuation line found in:\n    RUN exit. Empty continuation lines will become errors in a future release."}`,
	})
}

func Test400EmptyFile(t *testing.T) {
	makeTest(t, &param{
		expectedCode: 400,
		expectedBody: `{"error":true,"message":"Dockerfile may not be empty"}`,
	})
}

func Test400ParseInstructionError(t *testing.T) {
	makeTest(t, &param{
		body:         "FROM",
		expectedCode: 400,
		expectedBody: `{"error":true,"message":"Dockerfile parse error line 1: FROM requires either one or three arguments"}`,
	})
}

func Test200(t *testing.T) {
	makeTest(t, &param{
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

type param struct {
	path         string
	method       string
	body         string
	expectedCode int
	expectedBody string
}

func newParam(p *param) *param {
	if p.path == "" {
		p.path = "/"
	}
	if p.method == "" {
		p.method = "POST"
	}
	if p.expectedCode == 0 {
		p.expectedCode = 200
	}
	return p
}

func makeTest(t *testing.T, p *param) *httptest.ResponseRecorder {
	p = newParam(p)

	reqBody := bytes.NewReader([]byte(p.body))
	req, err := http.NewRequest(p.method, p.path, reqBody)
	equal(t, nil, err)

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(lintHandler)
	handler.ServeHTTP(rec, req)

	equal(t, p.expectedCode, rec.Code)
	equal(t, p.expectedBody, rec.Body.String())

	return rec
}
