package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)

	router.HandlerFunc(http.MethodGet, "/v1/:app/metadata/:type/:name", app.getMetadata)
	router.HandlerFunc(http.MethodGet, "/v1/:app/metadata-list/:type", app.getMetadataList)
	router.HandlerFunc(http.MethodPut, "/v1/:app/metadata/:type", app.putMetadata)
	router.HandlerFunc(http.MethodDelete, "/v1/:app/metadata/:type/:name", app.deleteMetadata)

	router.HandlerFunc(http.MethodGet, "/v1/attribute-types", app.getAllAttributeTypes)

	router.HandlerFunc(http.MethodGet, "/v1/:app/entity/:type/:name/:identifier", app.getOneEntity)
	router.HandlerFunc(http.MethodPost, "/v1/:app/search/:type/:entity/:index", app.searchEntities)
	router.HandlerFunc(http.MethodGet, "/v1/:app/list/:type/:name", app.getAllEntities)
	router.HandlerFunc(http.MethodPut, "/v1/:app/entity/:type/:name", app.putEntity)
	router.HandlerFunc(http.MethodDelete, "/v1/:app/entity/:type/:name/:identifier", app.deleteEntity)
	router.HandlerFunc(http.MethodGet, "/v1/types/:type/", app.getAllEntityTypes)

	router.HandlerFunc(http.MethodGet, "/v1/:app/edit/:type/:name/:identifier", app.editEntity)
	router.HandlerFunc(http.MethodGet, "/v1/:app/new/:type/:name", app.newEntity)
	router.HandlerFunc(http.MethodPost, "/v1/:app/post/:type/:name/:identifier", app.postEntity)
	router.HandlerFunc(http.MethodPost, "/v1/:app/post/:type/:name", app.postEntity)
	router.HandlerFunc(http.MethodGet, "/v1/l:app/ist-view/:type/:name", app.listEntities)

	// router.HandlerFunc(http.MethodGet, "/v1/app/edit/::app", app.editApp)
	// router.HandlerFunc(http.MethodGet, "/v1/app/new", app.newApp)
	// router.HandlerFunc(http.MethodPost, "/v1/app/post/:app", app.postApp)
	// router.HandlerFunc(http.MethodPost, "/v1/post/", app.postApp)
	// router.HandlerFunc(http.MethodGet, "/v1/app/list-view", app.listApps)

	router.HandlerFunc(http.MethodPost, "/v1/chatgpt/:type/:name", app.postChatGPT)

	return app.enableCORS(router)
}
