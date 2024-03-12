package main

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/ukubenet/metadata-repository/deployer"
	"github.com/ukubenet/metadata-repository/metadata"
	metaapi "github.com/ukubenet/metadata-repository/metadata/api"
)

func (app *application) getMetadata(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	name := params.ByName("name")
	entityType := params.ByName("type")

	entity, err := metaapi.ReadMetadata(metadata.EntityTypeMap[strings.ToLower(entityType)], name)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	parcel := getParcel(rw, r)
	parcel.Encode(http.StatusFound, entity)

}

func (app *application) getMetadataList(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	entityType := params.ByName("type")
	list, err := metaapi.ReadMetadataList(metadata.EntityTypeMap[strings.ToLower(entityType)])
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
	name := params.ByName("name")
	entityType := params.ByName("type")

	entity, err := metaapi.ReadMetadata(metadata.EntityTypeMap[strings.ToLower(entityType)], name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deployerFactory := deployer.CreateFactory(metadata.EntityTypeMap[strings.ToLower(entityType)])
	adapter := deployerFactory.CreateAdapter()
	err = adapter.Delete(entity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = metaapi.DeleteMetadata(metadata.EntityTypeMap[strings.ToLower(entityType)], name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) putMetadata(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	entityType := params.ByName("type")

	entityMetadata := new(metadata.EntityMetadata)

	parcel := getParcel(rw, r)
	err := parcel.Decode(entityMetadata)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = metaapi.PutMetadata(metadata.EntityTypeMap[strings.ToLower(entityType)], entityMetadata)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	deployerFactory := deployer.CreateFactory(metadata.EntityTypeMap[strings.ToLower(entityType)])
	adapter := deployerFactory.CreateAdapter()
	err = adapter.Deploy(entityMetadata)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}
