package main

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"net/http"
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
	app.logger.Println("id is:", id)

	entity, err := app.models.DB.Get(id)

	// just a sample how to generate a new UUID
	tempId := uuid.New()
	app.logger.Println("generated uuid is:", tempId.String())

	// todo: remove Entity sample.
	//entity := models.Entity{
	//	UUID:       tempId.String(),
	//	EntityName: "Some entity",
	//	CreatedAt:  time.Now(),
	//	UpdatedAt:  time.Now(),
	//}

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

}
