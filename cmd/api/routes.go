package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)

	router.HandlerFunc(http.MethodGet, "/v1/entity/:id", app.getOneEntity)
	router.HandlerFunc(http.MethodGet, "/v1/entities", app.getAllEntities)
	router.HandlerFunc(http.MethodPut, "/v1/entity", app.putEntity)

	return router
}
