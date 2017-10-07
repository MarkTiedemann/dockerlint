package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/docker/docker/builder/dockerfile/instructions"
	"github.com/docker/docker/builder/dockerfile/parser"
)

var (
	addr = flag.String("addr", ":3000", "the address of the server")
	path = flag.String("path", "/", "the path of the handler")
)

func main() {
	flag.Parse()
	http.HandleFunc("/", lintHandler)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func lintHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != *path {
		res.WriteHeader(404)
		return
	}

	if req.Method != "POST" {
		res.WriteHeader(405)
		res.Header().Set("Allow", "POST")
		return
	}

	ast, err := parser.Parse(req.Body)
	if err != nil {
		writeResult(res, 400, newLintError(err.Error()))
		return
	}

	if len(ast.Warnings) > 0 {
		writeResult(res, 400, newLintError(fmtWarnings(ast.Warnings)))
		return
	}

	if ast.AST.Children == nil {
		writeResult(res, 400, newLintError("Dockerfile may not be empty"))
		return
	}

	_, _, err = instructions.Parse(ast.AST)
	if err != nil {
		writeResult(res, 400, newLintError(err.Error()))
		return
	}

	writeResult(res, 200, &lintResult{Error: false})
}

type lintResult struct {
	Error   bool   `json:"error"`
	Message string `json:"message,omitempty"`
}

func newLintError(message string) *lintResult {
	return &lintResult{Error: true, Message: message}
}

func writeResult(res http.ResponseWriter, status int, result *lintResult) {
	res.WriteHeader(status)
	res.Header().Set("Content-Type", "application/json")
	buf, _ := json.Marshal(result)
	res.Write(buf)
}

func fmtWarnings(warnings []string) (fmted string) {
	slice := make([]string, 0)
	for _, warn := range warnings {
		slice = append(slice, strings.TrimPrefix(warn, "[WARNING]: "))
	}
	return strings.Join(slice, ". ")
}
