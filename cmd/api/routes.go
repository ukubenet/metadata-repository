package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)

	router.HandlerFunc(http.MethodGet, "/v1/entity/:name", app.getOneEntity)
	router.HandlerFunc(http.MethodGet, "/v1/entities", app.getAllEntities)
	router.HandlerFunc(http.MethodGet, "/v1/entities/:attribute_id", app.getAllEntitiesByAttribute)
	router.HandlerFunc(http.MethodPut, "/v1/entity", app.putEntity)
	router.HandlerFunc(http.MethodDelete, "/v1/entity/:name", app.deleteEntity)

	router.HandlerFunc(http.MethodGet, "/v1/attributes", app.getAllAttributes)

	return app.enableCORS(router)
}
