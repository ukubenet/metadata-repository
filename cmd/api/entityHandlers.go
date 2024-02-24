package main

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/ukubenet/metadata-repository/entity"
	entityapi "github.com/ukubenet/metadata-repository/entity/api"
	"github.com/ukubenet/metadata-repository/metadata"
)

func (app *application) getOneCatalogEntity(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	name := params.ByName("name")
	identifier := params.ByName("identifier")

	entity, err := entityapi.ReadEntity(metadata.Catalog, name, identifier)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	parcel := getParcel(rw, r)
	parcel.Encode(http.StatusFound, entity)
}

func (app *application) getAllCatalogEntities(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	name := params.ByName("name")

	list, err := entityapi.ReadEntities(metadata.Catalog, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parcel := getParcel(w, r)
	parcel.Encode(http.StatusOK, list)
}

func (app *application) putCatalogEntity(rw http.ResponseWriter, r *http.Request) {
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

func (app *application) deleteCatalogEntity(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	name := params.ByName("name")
	identifier := params.ByName("identifier")

	err := entityapi.DeleteEntity(metadata.Catalog, name, identifier)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (app *application) getAllCatalogEntityTypes(w http.ResponseWriter, r *http.Request) {
	list, err := entityapi.ReadEntityTypes(metadata.Catalog)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parcel := getParcel(w, r)
	parcel.Encode(http.StatusOK, list)
}

func (app *application) getOneEventEntity(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	name := params.ByName("name")
	identifier := params.ByName("identifier")

	entity, err := entityapi.ReadEntity(metadata.Event, name, identifier)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	parcel := getParcel(rw, r)
	parcel.Encode(http.StatusFound, entity)
}

func (app *application) getAllEventEntities(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	name := params.ByName("name")

	list, err := entityapi.ReadEntities(metadata.Event, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parcel := getParcel(w, r)
	parcel.Encode(http.StatusOK, list)
}

func (app *application) putEventEntity(rw http.ResponseWriter, r *http.Request) {
	entity := new(entity.EventEntity)

	parcel := getParcel(rw, r)
	err := parcel.Decode(entity)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if entity.EventTime.IsZero() {
		entity.EventTime = time.Now()
	}

	err = entityapi.PutEntity(entity)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (app *application) deleteEventEntity(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	name := params.ByName("name")
	identifier := params.ByName("identifier")

	err := entityapi.DeleteEntity(metadata.Event, name, identifier)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (app *application) getAllEventEntityTypes(w http.ResponseWriter, r *http.Request) {
	list, err := entityapi.ReadEntityTypes(metadata.Event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parcel := getParcel(w, r)
	parcel.Encode(http.StatusOK, list)
}
