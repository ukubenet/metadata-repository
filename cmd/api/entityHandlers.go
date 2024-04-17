package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/ukubenet/metadata-repository/entity"
	entityapi "github.com/ukubenet/metadata-repository/entity/api"
	indexItem "github.com/ukubenet/metadata-repository/entity/search/item"
	"github.com/ukubenet/metadata-repository/metadata"
)

func (app *application) getOneEntity(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	name := params.ByName("name")
	identifier := params.ByName("identifier")
	entityType := params.ByName("type")

	entType, ok := metadata.EntityTypeMap[strings.ToLower(entityType)]
	if !ok {
		http.Error(rw, "incorrect entity type: "+entityType, http.StatusBadRequest)
		return
	}

	entity, err := entityapi.ReadEntity(entType, name, identifier)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	parcel := getParcel(rw, r)
	parcel.Encode(http.StatusFound, entity)
}

func (app *application) searchEntities(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	entityName := params.ByName("entity")
	entityType := params.ByName("type")
	indexName := params.ByName("index")

	entType, ok := metadata.EntityTypeMap[strings.ToLower(entityType)]
	if !ok {
		http.Error(w, "incorrect entity type: "+entityType, http.StatusBadRequest)
		return
	}

	criteria := make(indexItem.ValueMap)

	parcel := getParcel(w, r)
	err := parcel.Decode(&criteria)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(criteria) == 0 {
		http.Error(w, "Criteria not defined", http.StatusBadRequest)
		return
	}

	list, err := entityapi.SearchEntities(entType, entityName, indexName, criteria)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parcel.Encode(http.StatusOK, list)
}

func (app *application) getAllEntities(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	name := params.ByName("name")
	entityType := params.ByName("type")

	entType, ok := metadata.EntityTypeMap[strings.ToLower(entityType)]
	if !ok {
		http.Error(w, "incorrect entity type: "+entityType, http.StatusBadRequest)
		return
	}

	list, err := entityapi.ReadEntities(entType, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parcel := getParcel(w, r)
	parcel.Encode(http.StatusOK, list)
}

func (app *application) putEntity(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	entityType := params.ByName("type")

	entType, ok := metadata.EntityTypeMap[strings.ToLower(entityType)]
	if !ok {
		http.Error(rw, "incorrect entity type: "+entityType, http.StatusBadRequest)
		return
	}
	e := entity.GetEntityInstance(entType)

	parcel := getParcel(rw, r)
	err := parcel.Decode(e)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	eventEntity, ok := e.(entity.EventEntity)
	if ok && eventEntity.EventTime.IsZero() {
		eventEntity.EventTime = time.Now()
	}

	err = entityapi.PutEntity(e)
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
	entityType := params.ByName("type")

	entType, ok := metadata.EntityTypeMap[strings.ToLower(entityType)]
	if !ok {
		http.Error(rw, "incorrect entity type: "+entityType, http.StatusBadRequest)
		return
	}

	err := entityapi.DeleteEntity(entType, name, identifier)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (app *application) getAllEntityTypes(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	entityType := params.ByName("type")

	entType, ok := metadata.EntityTypeMap[strings.ToLower(entityType)]
	if !ok {
		http.Error(w, "incorrect entity type: "+entityType, http.StatusBadRequest)
		return
	}

	list, err := entityapi.ReadEntityTypes(entType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parcel := getParcel(w, r)
	parcel.Encode(http.StatusOK, list)
}
