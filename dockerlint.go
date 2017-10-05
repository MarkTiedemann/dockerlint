package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/docker/docker/builder/dockerfile/instructions"
	"github.com/docker/docker/builder/dockerfile/parser"
)

type LintResult struct {
	Error   bool   `json:"error"`
	Message string `json:"message,omitempty"`
}

func lintError(message string) *LintResult {
	return &LintResult{Error: true, Message: message}
}

func writeJson(res http.ResponseWriter, status int, obj interface{}) {
	res.WriteHeader(status)
	res.Header().Set("Content-Type", "application/json")
	buf, _ := json.Marshal(obj)
	res.Write(buf)
}

func fmtWarnings(warnings []string) (fmted string) {
	slice := make([]string, 0)
	for _, warn := range warnings {
		slice = append(slice, strings.TrimPrefix(warn, "[WARNING]: "))
	}
	return strings.Join(slice, ". ")
}

func LintHandler(res http.ResponseWriter, req *http.Request) {
	ast, err := parser.Parse(req.Body)
	if err != nil {
		writeJson(res, 400, lintError(err.Error()))
		return
	}

	if len(ast.Warnings) > 0 {
		writeJson(res, 400, lintError(fmtWarnings(ast.Warnings)))
		return
	}

	if ast.AST.Children == nil {
		writeJson(res, 400, lintError("Dockerfile may not be empty"))
		return
	}

	_, _, err = instructions.Parse(ast.AST)
	if err != nil {
		writeJson(res, 400, lintError(err.Error()))
		return
	}

	writeJson(res, 200, &LintResult{Error: false})
}

func main() {
	http.HandleFunc("/", LintHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
