package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/skutaada/gonergy/lib"
	"github.com/skutaada/gonergy/router"
)

var (
	//go:embed all:templates/*
	templateFS embed.FS

	html *template.Template
)

func handleSignals(fetchC chan<- bool, serverC <-chan os.Signal, s *http.Server) {
	<-serverC
	if err := s.Shutdown(context.TODO()); err != nil {
		log.Fatal(err)
	}
	fetchC <- true
}

func main() {
	var err error
	html, err = router.TemplateParseFSRecursive(templateFS, ".html", true, nil)
	if err != nil {
		panic(err)
	}

	r := http.NewServeMux()
	r.Handle("GET /", router.Handler(index))

	server := http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	fetchC := make(chan bool)
	doneC := make(chan bool)
	serverC := make(chan os.Signal, 1)

	signal.Notify(serverC, syscall.SIGINT, syscall.SIGTERM)

	go lib.DailyFetchInsert2(fetchC, doneC)
	go handleSignals(fetchC, serverC, &server)

	fmt.Println("Listening on port 3000")

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
	<-doneC
	fmt.Println("Gracefully shutdown servers.")
}
