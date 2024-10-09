package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.ServeFiles("/static/*filepath",http.Dir("static"))

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)

	// Rest API - metadata
	router.HandlerFunc(http.MethodGet, "/v1/metadata/api/get/:app/:type/:name", app.getMetadata)
	router.HandlerFunc(http.MethodGet, "/v1/metadata/api/list/:app/:type", app.getMetadataList)
	router.HandlerFunc(http.MethodPut, "/v1/metadata/api/:app/:type", app.putMetadata)
	router.HandlerFunc(http.MethodDelete, "/v1/metadata/api/:app/:type/:name", app.deleteMetadata)

	router.HandlerFunc(http.MethodGet, "/v1/attribute-types", app.getAllAttributeTypes)

	// CRUD - metadata
	router.HandlerFunc(http.MethodGet, "/v1/metadata/edit/:app/:type/:name", app.editMetadata)
	router.HandlerFunc(http.MethodGet, "/v1/metadata/new/:app/:type", app.newMetadata)
	router.HandlerFunc(http.MethodPost, "/v1/metadata/chatgpt", app.postMetadataChatGPT)

	// Rest API - entity
	router.HandlerFunc(http.MethodGet, "/v1/entity/:app/:type/:name/:identifier", app.getOneEntity)
	router.HandlerFunc(http.MethodPost, "/v1/search/:app/:type/:entity/:index", app.searchEntities)
	router.HandlerFunc(http.MethodGet, "/v1/list/:app/:type/:name", app.getAllEntities)
	router.HandlerFunc(http.MethodPut, "/v1/entity/:app/:type/:name", app.putEntity)
	router.HandlerFunc(http.MethodDelete, "/v1/entity/:app/:type/:name/:identifier", app.deleteEntity)

	router.HandlerFunc(http.MethodGet, "/v1/types/:type/", app.getAllEntityTypes)

	// CRUD - entity
	router.HandlerFunc(http.MethodGet, "/v1/edit/:app/:type/:name/:identifier", app.editEntity)
	router.HandlerFunc(http.MethodGet, "/v1/new/:app/:type/:name", app.newEntity)
	router.HandlerFunc(http.MethodPost, "/v1/post/:app/:type/:name/:identifier", app.postEntity)
	router.HandlerFunc(http.MethodPost, "/v1/post/:app/:type/:name", app.postEntity)
	router.HandlerFunc(http.MethodGet, "/v1/list-view/:app/:type/:name", app.listEntities)
	router.HandlerFunc(http.MethodPost, "/v1/chatgpt/:app/:type/:name", app.postEntityChatGPT)

	// CRUD - app
	router.HandlerFunc(http.MethodGet, "/v1/app/edit/:app", app.editApp)
	router.HandlerFunc(http.MethodGet, "/v1/app/duplicate/:app", app.duplicateApp)
	router.HandlerFunc(http.MethodGet, "/v1/app/new", app.newApp)
	router.HandlerFunc(http.MethodPost, "/v1/app/post/:app", app.postApp)
	router.HandlerFunc(http.MethodPost, "/v1/app/post", app.postNewApp)
	router.HandlerFunc(http.MethodDelete, "/v1/app/delete/:app", app.deleteApp)
	router.HandlerFunc(http.MethodGet, "/v1/app/config/:app", app.appConfiguration)
	router.HandlerFunc(http.MethodGet, "/v1/app/run/:app", app.appRun)
	router.HandlerFunc(http.MethodGet, "/v1/", app.listApps)
	router.HandlerFunc(http.MethodPost, "/v1/app/chatgpt", app.postAppChatGPT)

	return app.enableCORS(router)
}
