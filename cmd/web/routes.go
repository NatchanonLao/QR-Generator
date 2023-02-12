package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	fileServer := http.FileServer(http.Dir("./static/"))
	router.Use(app.rateLimit)
	router.Get("/healthcheck", app.healthcheckHandler)
	router.Get("/qr", app.homeHandler)
	router.Post("/qr", app.generateQRHandler)
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return router
}
