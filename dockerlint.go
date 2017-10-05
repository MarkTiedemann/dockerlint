package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/docker/docker/builder/dockerfile/instructions"
	"github.com/docker/docker/builder/dockerfile/parser"
)

type Error struct {
	Error    bool     `json:"error"`
	Messages []string `json:"messages"`
}

func writeJson(res http.ResponseWriter, status int, obj interface{}) {
	res.WriteHeader(status)
	res.Header().Set("Content-Type", "application/json")
	buf, _ := json.Marshal(obj)
	res.Write(buf)
}

func formatWarnings(warnings []string) (fmted []string) {
	for _, warn := range warnings {
		fmted = append(fmted, strings.TrimPrefix(warn, "[WARNING]: "))
	}
	return fmted
}

func LintHandler(res http.ResponseWriter, req *http.Request) {
	ast, err := parser.Parse(req.Body)
	if err != nil {
		writeJson(res, 400, &Error{Error: true, Messages: []string{err.Error()}})
		return
	}

	if len(ast.Warnings) > 0 {
		writeJson(res, 400, &Error{Error: true, Messages: formatWarnings(ast.Warnings)})
		return
	}

	if ast.AST.Children == nil {
		writeJson(res, 400, &Error{Error: true, Messages: []string{"Dockerfile may not be empty"}})
		return
	}

	_, _, err = instructions.Parse(ast.AST)
	if err != nil {
		writeJson(res, 400, &Error{Error: true, Messages: []string{err.Error()}})
		return
	}

	writeJson(res, 200, &Error{Error: false, Messages: []string{}})
}

func main() {
	http.HandleFunc("/", LintHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
