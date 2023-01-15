package main

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"net/http"

	"github.com/ukubenet/metadata-repository/models"
	parcel "github.com/ukubenet/metadata-repository/parser"
)

func (app *application) getOneEntity(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	// use commented code if "id" is int type
	//id, err := strconv.Atoi(params.ByName("id"))
	//if err != nil {
	//	app.logger.Print(errors.New("invalid id parameter"))
	//	app.errorJSON(w, err)
	//	return
	//}
	id := params.ByName("id")

	entity, err := app.models.DB.Get(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// just a sample how to generate a new UUID
	tempId := uuid.New()
	app.logger.Println("generated uuid is:", tempId.String())

	err = app.writeJSON(w, http.StatusOK, entity, "entity")
	if err != nil {
		app.logger.Print(errors.New("invalid id parameter"))
		app.errorJSON(w, err)
		return
	}

	js, err := json.Marshal(entity)
	if err != nil {
		app.logger.Println("error write to Redis:", err)
	}
	app.writeRedis("entity:"+entity.UUID, string(js))
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
	factory := parcel.CreateFactory()
	p := factory.Parcel(rw, r)

	entity := new(models.Entity)
	p.Decode(entity)
	p.Encode(http.StatusCreated, entity)
}
