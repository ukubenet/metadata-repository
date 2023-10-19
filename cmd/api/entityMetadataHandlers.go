package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ukubenet/metadata-repository/metadata"
	metaapi "github.com/ukubenet/metadata-repository/metadata/api"
)

func (app *application) getOneEntityMetadata(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	name := params.ByName("name")

	entity, err := metaapi.ReadMetadata(name)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	parcel := getParcel(rw, r)
	parcel.Encode(http.StatusFound, entity)

}

func (app *application) getAllEntityMetadataList(w http.ResponseWriter, r *http.Request) {
	list, err := metaapi.ReadMetadataList()
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

func (app *application) deleteEntityMetadata(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	name := params.ByName("name")

	err := metaapi.DeleteMetadata(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) putEntityMetadata(rw http.ResponseWriter, r *http.Request) {
	metadata := new(metadata.EntityMetadata)

	parcel := getParcel(rw, r)
	err := parcel.Decode(metadata)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = metaapi.PutMetadata(metadata)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}
