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
	adapter.Read(name, entity)

	output := parcel.CreateFactory()
	parcel := output.Parcel(rw, r)

	parcel.Encode(http.StatusFound, entity)

}

func (app *application) getAllEntities(w http.ResponseWriter, r *http.Request) {
	storage := storage.CreateFactory()
	adapter := storage.CreateAdapter()
	list := []string{}
	adapter.List(&list)

	output := parcel.CreateFactory()
	parcel := output.Parcel(w, r)

	parcel.Encode(http.StatusOK, list)
}

func (app *application) getAllAttributes(w http.ResponseWriter, r *http.Request) {

}

func (app *application) deleteEntity(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	name := params.ByName("name")

	storage := storage.CreateFactory()
	adapter := storage.CreateAdapter()
	adapter.Delete(name)

	w.WriteHeader(http.StatusOK)
}

func (app *application) putEntity(rw http.ResponseWriter, r *http.Request) {
	entity := new(models.Entity)

	factoryReader := parcel.CreateFactory()
	parcelReader := factoryReader.Parcel(rw, r)
	parcelReader.Decode(entity)

	factoryWriter := storage.CreateFactory()
	adapter := factoryWriter.CreateAdapter()
	adapter.Put(entity)

	rw.WriteHeader(http.StatusCreated)
}
