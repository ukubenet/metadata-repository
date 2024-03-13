package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)

	router.HandlerFunc(http.MethodGet, "/v1/metadata/:type/:name", app.getMetadata)
	router.HandlerFunc(http.MethodGet, "/v1/metadata-list/:type", app.getMetadataList)
	router.HandlerFunc(http.MethodPut, "/v1/metadata/:type", app.putMetadata)
	router.HandlerFunc(http.MethodDelete, "/v1/metadata/:type/:name", app.deleteMetadata)

	router.HandlerFunc(http.MethodGet, "/v1/attribute-types", app.getAllAttributeTypes)

	router.HandlerFunc(http.MethodGet, "/v1/entity/:type/:name/:identifier", app.getOneEntity)
	router.HandlerFunc(http.MethodGet, "/v1/list/:type/:name", app.getAllEntities)
	router.HandlerFunc(http.MethodPut, "/v1/entity/:type/:name", app.putEntity)
	router.HandlerFunc(http.MethodDelete, "/v1/entity/:type/:name/:identifier", app.deleteEntity)
	router.HandlerFunc(http.MethodGet, "/v1/:type/", app.getAllEntityTypes)

	return app.enableCORS(router)
}
