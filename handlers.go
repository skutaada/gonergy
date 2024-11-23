package main

import (
	"net/http"

	"github.com/skutaada/gonergy/router"
)

func index(_ *http.Request) *router.Response {
	return router.HTML(http.StatusOK, html, "index.html", nil, nil)
}
