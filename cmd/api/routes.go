package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)

	router.HandlerFunc(http.MethodGet, "/v1/entity-metadata/:name", app.getOneEntityMetadata)
	router.HandlerFunc(http.MethodGet, "/v1/entity-metadata-list", app.getAllEntityMetadataList)
	router.HandlerFunc(http.MethodPut, "/v1/entity-metadata", app.putEntityMetadata)
	router.HandlerFunc(http.MethodDelete, "/v1/entity-metadata/:name", app.deleteEntityMetadata)

	router.HandlerFunc(http.MethodGet, "/v1/attribute-types", app.getAllAttributeTypes)

	router.HandlerFunc(http.MethodGet, "/v1/entity/:name/:identifier", app.getOneEntity)
	router.HandlerFunc(http.MethodGet, "/v1/entities/:name", app.getAllEntities)
	router.HandlerFunc(http.MethodPut, "/v1/entity/:name", app.putEntity)
	router.HandlerFunc(http.MethodDelete, "/v1/entity/:name/:identifier", app.deleteEntity)
	router.HandlerFunc(http.MethodGet, "/v1/entities", app.getAllEntityTypes)

	return app.enableCORS(router)
}
