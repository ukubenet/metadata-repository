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

	router.HandlerFunc(http.MethodGet, "/v1/event-metadata/:name", app.getCatalogMetadata)
	router.HandlerFunc(http.MethodGet, "/v1/event-metadata-list", app.getCatalogMetadataList)
	router.HandlerFunc(http.MethodPut, "/v1/event-metadata", app.putCatalogMetadata)
	router.HandlerFunc(http.MethodDelete, "/v1/event-metadata/:name", app.deleteCatalogMetadata)

	router.HandlerFunc(http.MethodGet, "/v1/attribute-types", app.getAllAttributeTypes)

	router.HandlerFunc(http.MethodGet, "/v1/catalog-entity/:name/:identifier", app.getOneEntity)
	router.HandlerFunc(http.MethodGet, "/v1/catalog-list/:name", app.getAllEntities)
	router.HandlerFunc(http.MethodPut, "/v1/catalog-entity/:name", app.putEntity)
	router.HandlerFunc(http.MethodDelete, "/v1/catalog-entity/:name/:identifier", app.deleteEntity)
	router.HandlerFunc(http.MethodGet, "/v1/catalogs", app.getAllEntityTypes)

	return app.enableCORS(router)
}
