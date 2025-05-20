package main

import (
	"net/http"
	"net/url"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/ukubenet/metadata-repository/entity"
	entityapi "github.com/ukubenet/metadata-repository/entity/api"
	global "github.com/ukubenet/metadata-repository/global"
)

func (app *application) readEventRegisters(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")
	global.SetAppName(appName)

	event := params.ByName("event")
	identifier := params.ByName("identifier")

	registers, err := entityapi.ReadEventRegisters(event, identifier)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	parcel := getParcel(rw, r)
	parcel.Encode(http.StatusOK, registers)
}

func (app *application) commitEvent(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")
	global.SetAppName(appName)

	event := params.ByName("event")
	identifier := params.ByName("identifier")

	registers := new([]*entity.Register)

	parcel := getParcel(rw, r)
	err := parcel.Decode(registers)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = entityapi.CommitEvent(event, identifier, *registers)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (app *application) rollbackEvent(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")
	global.SetAppName(appName)

	event := params.ByName("event")
	identifier := params.ByName("identifier")

	err := entityapi.RollbackEvent(event, identifier)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (app *application) readState(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	appName := params.ByName("app")
	global.SetAppName(appName)

	registerType := params.ByName("type")
	registerName := params.ByName("register")

	// regType, ok := metadata.RegisterTypeMap[strings.ToLower(registerType)]
	// if !ok {
	// 	http.Error(w, "incorrect entity type: "+registerType, http.StatusBadRequest)
	// 	return
	// }

	timestamp, err := url.PathUnescape(params.ByName("timestamp"))
	if err != nil {
		http.Error(w, "Invalid register name", http.StatusBadRequest)
		return
	}

	dimensions := new([]entity.AttributeValues)
	parcel := getParcel(w, r)
	err = parcel.Decode(dimensions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, timestamp)
	if err != nil {
		http.Error(w, "Invalid register name", http.StatusBadRequest)
		return
	}

	state, err := entityapi.ReadState(registerType, registerName, *dimensions, t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parcel.Encode(http.StatusOK, state)
}
