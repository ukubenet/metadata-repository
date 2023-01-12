package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	parcel "github.com/ukubenet/metadata-repository/parser"
	"github.com/ukubenet/metadata-repository/parser/encoding"
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

func (app *application) deleteEntity(w http.ResponseWriter, r *http.Request) {

}

func (app *application) insertEntity(w http.ResponseWriter, r *http.Request) {

}

func (app *application) updateEntity(w http.ResponseWriter, r *http.Request) {

}

func (app *application) searchEntities(w http.ResponseWriter, r *http.Request) {

}

func (app *application) putEntity(rw http.ResponseWriter, r *http.Request) {
	// New factory
	factory := parcel.NewFactory()

	// Encoders/Decoders will be called in the order
	// they are registered. The setup below:
	// Request ->
	// 1. Query Strings
	// 2. Json
	// 3. Xml
	// Response ->
	// 1. Json
	// 2. Xml
	// Notice that the Query codec only provides a decoder,
	// so it will not be added to response chain
	factory.Use(encoding.Query())
	factory.Use(encoding.JSON())
	factory.Use(encoding.XML())

	p := factory.Parcel(rw, r)
	type (
		TestPerson struct {
			Name         string  `xml:"name" json:"name"`
			Email        string  `xml:"email" json:"email"`
			IsAdmin      bool    `xml:"is-admin" json:"isAdmin"`
			Age          int     `xml:"age" json:"age"`
			HourlyRate   float32 `xml:"hourly-rate" json:"hourlyRate"`
			AccessToken  string  `query:"access-token"`
			ShowMatching bool    `query:"show-matching"`
		}
	)
	person1 := new(TestPerson)

	p.Decode(person1)

	p.Encode(http.StatusCreated, person1)
}
