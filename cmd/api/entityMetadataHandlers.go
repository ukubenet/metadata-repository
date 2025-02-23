package main

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/ukubenet/metadata-repository/deployer"
	global "github.com/ukubenet/metadata-repository/global"
	"github.com/ukubenet/metadata-repository/metadata"
	metaapi "github.com/ukubenet/metadata-repository/metadata/api"
)

func (app *application) getMetadata(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")
	global.SetAppName(appName)

	name := params.ByName("name")
	entityType := params.ByName("type")

	entType, ok := metadata.EntityTypeMap[strings.ToLower(entityType)]
	if !ok {
		http.Error(rw, "incorrect entity type: "+entityType, http.StatusBadRequest)
		return
	}

	entity, err := metaapi.ReadMetadata(entType, name)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	parcel := getParcel(rw, r)
	parcel.Encode(http.StatusFound, entity)

}

func (app *application) getMetadataList(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")
	global.SetAppName(appName)

	entityType := params.ByName("type")

	entType, ok := metadata.EntityTypeMap[strings.ToLower(entityType)]
	if !ok {
		http.Error(w, "incorrect entity type: "+entityType, http.StatusBadRequest)
		return
	}

	list, err := metaapi.ReadMetadataList(entType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parcel := getParcel(w, r)
	parcel.Encode(http.StatusOK, list)
}

func (app *application) getAllAttributeTypes(w http.ResponseWriter, r *http.Request) {
	parcel := getParcel(w, r)
	parcel.Encode(http.StatusOK, metadata.AttributeTypeList)
}

func (app *application) deleteMetadata(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")
	global.SetAppName(appName)

	name := params.ByName("name")
	entityType := params.ByName("type")

	entType, ok := metadata.EntityTypeMap[strings.ToLower(entityType)]
	if !ok {
		http.Error(w, "incorrect entity type: "+entityType, http.StatusBadRequest)
		return
	}

	entity, err := metaapi.ReadMetadata(entType, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deployerFactory := deployer.CreateFactory(entType)
	adapter := deployerFactory.CreateAdapter()
	err = adapter.Delete(entity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = metaapi.DeleteMetadata(entType, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) putMetadata(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")
	global.SetAppName(appName)

	entityType := params.ByName("type")

	entityMetadata := new(metadata.EntityMetadata)

	parcel := getParcel(rw, r)
	err := parcel.Decode(entityMetadata)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	entType, ok := metadata.EntityTypeMap[strings.ToLower(entityType)]
	if !ok {
		http.Error(rw, "incorrect entity type: "+entityType, http.StatusBadRequest)
		return
	}

	if entType == metadata.Event {
		if _, ok := entityMetadata.Attributes["EventTime"]; !ok {
			entityMetadata.Attributes["EventTime"] = metadata.Attribute{
				"type": "dateTime",
			}
		}
	}

	err = metaapi.PutMetadata(entType, entityMetadata)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// if entityType == "event" {
	// 	path, _ := os.Getwd()
	// 	tmplFile := path + "/" + config.Config.Metadata.Path + "/" + global.AppName + "/" + config.Config.Metadata.Tmplsubpath + "/" + strings.ToLower(entityMetadata.EntityName) + "_transaction.tmpl"
	// 	err = ioutil.WriteFile(tmplFile, []byte(entityMetadata.CustomTemplate), 0644)
	// 	if err != nil {
	// 		http.Error(rw, err.Error(), http.StatusBadRequest)
	// 	}
	// }

	deployerFactory := deployer.CreateFactory(entType)
	adapter := deployerFactory.CreateAdapter()
	err = adapter.Deploy(entityMetadata)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (app *application) copyMetadata(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")
	global.SetAppName(appName)

	name := params.ByName("name")
	entityType := params.ByName("type")

	entType, ok := metadata.EntityTypeMap[strings.ToLower(entityType)]
	if !ok {
		http.Error(rw, "incorrect entity type: "+entityType, http.StatusBadRequest)
		return
	}

	entityMetadata, err := metaapi.ReadMetadata(entType, name)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	entityMetadata.EntityName = entityMetadata.EntityName + "_copy"

	err = metaapi.PutMetadata(entType, entityMetadata)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// if entityType == "event" {
	// 	path, _ := os.Getwd()
	// 	tmplFile := path + "/" + config.Config.Metadata.Path + "/" + global.AppName + "/" + config.Config.Metadata.Tmplsubpath + "/" + strings.ToLower(entityMetadata.EntityName) + "_transaction.tmpl"
	// 	err = ioutil.WriteFile(tmplFile, []byte(entityMetadata.CustomTemplate), 0644)
	// 	if err != nil {
	// 		http.Error(rw, err.Error(), http.StatusBadRequest)
	// 	}
	// }

	deployerFactory := deployer.CreateFactory(entType)
	adapter := deployerFactory.CreateAdapter()
	err = adapter.Deploy(entityMetadata)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusCreated)

}
