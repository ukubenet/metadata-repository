package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/ukubenet/metadata-repository/models"
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

	tempId := uuid.New()
	app.logger.Println("id is:", id)
	app.logger.Println("uuid is:", tempId.String())

	entity := models.Entity{
		UUID:       tempId.String(),
		EntityName: "Some entity",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := app.writeJSON(w, http.StatusOK, entity, "entity")
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
