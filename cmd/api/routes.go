package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)

	router.HandlerFunc(http.MethodGet, "/v1/catalog-metadata/:name", app.getCatalogMetadata)
	router.HandlerFunc(http.MethodGet, "/v1/catalog-metadata-list", app.getCatalogMetadataList)
	router.HandlerFunc(http.MethodPut, "/v1/catalog-metadata", app.putCatalogMetadata)
	router.HandlerFunc(http.MethodDelete, "/v1/catalog-metadata/:name", app.deleteCatalogMetadata)

	router.HandlerFunc(http.MethodGet, "/v1/event-metadata/:name", app.getEventMetadata)
	router.HandlerFunc(http.MethodGet, "/v1/event-metadata-list", app.getEventMetadataList)
	router.HandlerFunc(http.MethodPut, "/v1/event-metadata", app.putEventMetadata)
	router.HandlerFunc(http.MethodDelete, "/v1/event-metadata/:name", app.deleteEventMetadata)

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
