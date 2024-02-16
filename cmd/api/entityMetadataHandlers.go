package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ukubenet/metadata-repository/deployer"
	"github.com/ukubenet/metadata-repository/metadata"
	metaapi "github.com/ukubenet/metadata-repository/metadata/api"
)

func (app *application) getCatalogMetadata(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	name := params.ByName("name")

	entity, err := metaapi.ReadCatalogMetadata(name)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	parcel := getParcel(rw, r)
	parcel.Encode(http.StatusFound, entity)

}

func (app *application) getCatalogMetadataList(w http.ResponseWriter, r *http.Request) {
	list, err := metaapi.ReadCatalogMetadataList()
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

func (app *application) deleteCatalogMetadata(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	name := params.ByName("name")

	entity, err := metaapi.ReadCatalogMetadata(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deployerFactory := deployer.CreateFactory(metadata.Catalog)
	adapter := deployerFactory.CreateAdapter()
	err = adapter.Delete(entity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = metaapi.DeleteCatalogMetadata(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) putCatalogMetadata(rw http.ResponseWriter, r *http.Request) {
	entityMetadata := new(metadata.EntityMetadata)

	parcel := getParcel(rw, r)
	err := parcel.Decode(entityMetadata)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = metaapi.PutCatalogMetadata(entityMetadata)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	deployerFactory := deployer.CreateFactory(metadata.Catalog)
	adapter := deployerFactory.CreateAdapter()
	err = adapter.Deploy(entityMetadata)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}
