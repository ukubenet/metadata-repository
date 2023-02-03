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
	entities, err := app.models.DB.All()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, entities, "entities")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) getAllAttributes(w http.ResponseWriter, r *http.Request) {
	attributes, err := app.models.DB.AttributesAll()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, attributes, "attributes")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) getAllEntitiesByAttribute(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	attributeID := params.ByName("attribute_id")

	entities, err := app.models.DB.All(attributeID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, entities, "entities")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) deleteEntity(w http.ResponseWriter, r *http.Request) {

}

func (app *application) insertEntity(w http.ResponseWriter, r *http.Request) {

}

func (app *application) updateEntity(w http.ResponseWriter, r *http.Request) {

}

func (app *application) searchEntities(w http.ResponseWriter, r *http.Request) {

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
