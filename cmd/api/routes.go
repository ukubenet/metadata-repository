package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)

	router.HandlerFunc(http.MethodGet, "/v1/metadata/:app/:type/:name", app.getMetadata)
	router.HandlerFunc(http.MethodGet, "/v1/metadata-list/:app/:type", app.getMetadataList)
	router.HandlerFunc(http.MethodPut, "/v1/metadata/:app/:type", app.putMetadata)
	router.HandlerFunc(http.MethodDelete, "/v1/metadata/:app/:type/:name", app.deleteMetadata)

	router.HandlerFunc(http.MethodGet, "/v1/attribute-types", app.getAllAttributeTypes)

	router.HandlerFunc(http.MethodGet, "/v1/entity/:app/:type/:name/:identifier", app.getOneEntity)
	router.HandlerFunc(http.MethodPost, "/v1/search/:app/:type/:entity/:index", app.searchEntities)
	router.HandlerFunc(http.MethodGet, "/v1/list/:app/:type/:name", app.getAllEntities)
	router.HandlerFunc(http.MethodPut, "/v1/entity/:app/:type/:name", app.putEntity)
	router.HandlerFunc(http.MethodDelete, "/v1/entity/:app/:type/:name/:identifier", app.deleteEntity)

	router.HandlerFunc(http.MethodGet, "/v1/types/:type/", app.getAllEntityTypes)

	router.HandlerFunc(http.MethodGet, "/v1/edit/:app/:type/:name/:identifier", app.editEntity)
	router.HandlerFunc(http.MethodGet, "/v1/new/:app/:type/:name", app.newEntity)
	router.HandlerFunc(http.MethodPost, "/v1/post/:app/:type/:name/:identifier", app.postEntity)
	router.HandlerFunc(http.MethodPost, "/v1/post/:app/:type/:name", app.postEntity)
	router.HandlerFunc(http.MethodGet, "/v1/list-view/:app/:type/:name", app.listEntities)
	router.HandlerFunc(http.MethodPost, "/v1/chatgpt/:app/:type/:name", app.postEntityChatGPT)

	router.HandlerFunc(http.MethodGet, "/v1/app/edit/:app", app.editApp)
	router.HandlerFunc(http.MethodGet, "/v1/app/new", app.newApp)
	router.HandlerFunc(http.MethodPost, "/v1/app/post/:app", app.postApp)
	router.HandlerFunc(http.MethodPost, "/v1/app/post", app.postNewApp)
	router.HandlerFunc(http.MethodDelete, "/v1/app/delete/:app", app.deleteApp)
	router.HandlerFunc(http.MethodGet, "/v1/app/config/:app", app.appConfiguration)
	// router.HandlerFunc(http.MethodPost, "/v1/app/:app", app.runApp)
	// router.HandlerFunc(http.MethodPost, "/v1/app/run/:app", app.runApp)
	router.HandlerFunc(http.MethodGet, "/v1/", app.listApps)
	router.HandlerFunc(http.MethodPost, "/v1/app/chatgpt", app.postAppChatGPT)

	return app.enableCORS(router)
}
