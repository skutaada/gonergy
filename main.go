package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/skutaada/gonergy/router"
)

var (
	//go:embed all:templates/*
	templateFS embed.FS

	html *template.Template
)

func main() {
	var err error
	html, err = router.TemplateParseFSRecursive(templateFS, ".html", true, nil)
	if err != nil {
		panic(err)
	}

	r := http.NewServeMux()
	r.Handle("GET /", router.Handler(index))

	fmt.Println("Listening on port 3000")

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
