package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ukubenet/metadata-repository/entity"
	entityapi "github.com/ukubenet/metadata-repository/entity/api"
)

func (app *application) getOneEntity(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	name := params.ByName("name")
	identifier := params.ByName("identifier")

	entity, err := entityapi.ReadEntity(name, identifier)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	parcel := getParcel(rw, r)
	parcel.Encode(http.StatusFound, entity)
}

func (app *application) getAllEntities(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	name := params.ByName("name")

	list, err := entityapi.ReadEntities(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parcel := getParcel(w, r)
	parcel.Encode(http.StatusOK, list)
}

func (app *application) putEntity(rw http.ResponseWriter, r *http.Request) {
	entity := new(entity.CatalogEntity)

	parcel := getParcel(rw, r)
	err := parcel.Decode(entity)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = entityapi.PutEntity(entity)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (app *application) deleteEntity(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	name := params.ByName("name")
	identifier := params.ByName("identifier")

	err := entityapi.DeleteEntity(name, identifier)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (app *application) getAllEntityTypes(w http.ResponseWriter, r *http.Request) {
	list, err := entityapi.ReadEntityTypes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parcel := getParcel(w, r)
	parcel.Encode(http.StatusOK, list)
}
