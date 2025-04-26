package main

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/ukubenet/metadata-repository/deployer"
	global "github.com/ukubenet/metadata-repository/global"
	"github.com/ukubenet/metadata-repository/metadata"
	metaapi "github.com/ukubenet/metadata-repository/metadata/api"
)

func (app *application) getRegisterMetadata(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")
	global.SetAppName(appName)

	name := params.ByName("name")
	registerType := params.ByName("type")

	regType, ok := metadata.RegisterTypeMap[strings.ToLower(registerType)]
	if !ok {
		http.Error(rw, "incorrect entity type: "+registerType, http.StatusBadRequest)
		return
	}

	register, err := metaapi.ReadRegisterMetadata(regType, name)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	parcel := getParcel(rw, r)
	parcel.Encode(http.StatusOK, register)

}

func (app *application) getRegisterMetadataList(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")
	global.SetAppName(appName)

	registerType := params.ByName("type")

	regType, ok := metadata.RegisterTypeMap[strings.ToLower(registerType)]
	if !ok {
		http.Error(w, "incorrect entity type: "+registerType, http.StatusBadRequest)
		return
	}

	list, err := metaapi.ReadRegisterMetadataList(regType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parcel := getParcel(w, r)
	parcel.Encode(http.StatusOK, list)
}

func (app *application) deleteRegisterMetadata(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")
	global.SetAppName(appName)

	name := params.ByName("name")
	registerType := params.ByName("type")

	regType, ok := metadata.RegisterTypeMap[strings.ToLower(registerType)]
	if !ok {
		http.Error(w, "incorrect entity type: "+registerType, http.StatusBadRequest)
		return
	}

	register, err := metaapi.ReadRegisterMetadata(regType, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deployerFactory := deployer.CreateRegisterFactory(regType)
	adapter := deployerFactory.CreateAdapter()
	err = adapter.RegisterDelete(register)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = metaapi.DeleteRegisterMetadata(regType, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) putRegisterMetadata(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")
	global.SetAppName(appName)

	registerType := params.ByName("type")

	registerMetadata := new(metadata.RegisterMetadata)

	parcel := getParcel(rw, r)
	err := parcel.Decode(registerMetadata)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	regType, ok := metadata.RegisterTypeMap[strings.ToLower(registerType)]
	if !ok {
		http.Error(rw, "incorrect entity type: "+registerType, http.StatusBadRequest)
		return
	}

	err = metaapi.PutRegisterMetadata(regType, registerMetadata)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	deployerFactory := deployer.CreateRegisterFactory(regType)
	adapter := deployerFactory.CreateAdapter()
	err = adapter.RegisterDeploy(registerMetadata)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (app *application) copyRegisterMetadata(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")
	global.SetAppName(appName)

	name := params.ByName("name")
	registerType := params.ByName("type")

	regType, ok := metadata.RegisterTypeMap[strings.ToLower(registerType)]
	if !ok {
		http.Error(rw, "incorrect entity type: "+registerType, http.StatusBadRequest)
		return
	}

	registerMetadata, err := metaapi.ReadRegisterMetadata(regType, name)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	registerMetadata.RegisterName = registerMetadata.RegisterName + "_copy"

	err = metaapi.PutRegisterMetadata(regType, registerMetadata)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	deployerFactory := deployer.CreateRegisterFactory(regType)
	adapter := deployerFactory.CreateAdapter()
	err = adapter.RegisterDeploy(registerMetadata)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}
