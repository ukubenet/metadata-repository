package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/ukubenet/metadata-repository/models"
	parcel "github.com/ukubenet/metadata-repository/parser"
	storage "github.com/ukubenet/metadata-repository/storage"
)

func (app *application) getOneEntity(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	entity := new(models.Entity)
	name := params.ByName("name")

	dbReader := storage.CreateFactory()
	adapter := dbReader.CreateAdapter()
	err := adapter.Read(name, entity)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	output := parcel.CreateFactory()
	parcel := output.Parcel(rw, r)

	parcel.Encode(http.StatusFound, entity)

}

func (app *application) getAllEntities(w http.ResponseWriter, r *http.Request) {
	storage := storage.CreateFactory()
	adapter := storage.CreateAdapter()
	list := []string{}
	err := adapter.List(&list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output := parcel.CreateFactory()
	parcel := output.Parcel(w, r)

	parcel.Encode(http.StatusOK, list)
}

func (app *application) getAllAttributeTypes(w http.ResponseWriter, r *http.Request) {
	list := []string{"string", "integer", "decimal", "datetime", "boolean"}

	output := parcel.CreateFactory()
	parcel := output.Parcel(w, r)

	parcel.Encode(http.StatusOK, list)
}

func (app *application) deleteEntity(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	name := params.ByName("name")

	storage := storage.CreateFactory()
	adapter := storage.CreateAdapter()
	err := adapter.Delete(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) putEntity(rw http.ResponseWriter, r *http.Request) {
	entity := new(models.Entity)

	factoryReader := parcel.CreateFactory()
	parcelReader := factoryReader.Parcel(rw, r)
	err := parcelReader.Decode(entity)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if entity.EntityName == "" {
		http.Error(rw, "Entity name not defined", http.StatusBadRequest)
		return
	}
	if len(entity.Attributes) == 0 {
		http.Error(rw, "Entity attributes not defined", http.StatusBadRequest)
		return
	}

	factoryWriter := storage.CreateFactory()
	adapter := factoryWriter.CreateAdapter()
	err = adapter.Put(entity)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}
