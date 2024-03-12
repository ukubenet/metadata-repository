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

	router.HandlerFunc(http.MethodGet, "/v1/catalog-entity/:name/:identifier", app.getOneCatalogEntity)
	router.HandlerFunc(http.MethodGet, "/v1/catalog-list/:name", app.getAllCatalogEntities)
	router.HandlerFunc(http.MethodPut, "/v1/catalog-entity/:name", app.putCatalogEntity)
	router.HandlerFunc(http.MethodDelete, "/v1/catalog-entity/:name/:identifier", app.deleteCatalogEntity)
	router.HandlerFunc(http.MethodGet, "/v1/catalogs", app.getAllCatalogEntityTypes)

	router.HandlerFunc(http.MethodGet, "/v1/event-entity/:name/:identifier", app.getOneEventEntity)
	router.HandlerFunc(http.MethodGet, "/v1/event-list/:name", app.getAllEventEntities)
	router.HandlerFunc(http.MethodPut, "/v1/event-entity/:name", app.putEventEntity)
	router.HandlerFunc(http.MethodDelete, "/v1/event-entity/:name/:identifier", app.deleteEventEntity)
	router.HandlerFunc(http.MethodGet, "/v1/events", app.getAllEventEntityTypes)

	return app.enableCORS(router)
}
