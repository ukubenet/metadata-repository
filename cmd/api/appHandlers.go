package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/ukubenet/metadata-repository/config"
	global "github.com/ukubenet/metadata-repository/global"
	"github.com/ukubenet/metadata-repository/metadata"
	metaapi "github.com/ukubenet/metadata-repository/metadata/api"
)

func (app *application) apiListApps(w http.ResponseWriter, r *http.Request) {
	// todo: should be adapter to get app list from metadata repository
	list := []string{}
	path, _ := os.Getwd()
	appPath := path + "/" + config.Config.Metadata.Path
	files, err := os.ReadDir(appPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			if file.Name() == ".git" || file.Name() == "metadata" {
				continue
			}
			list = append(list, file.Name())
		}
	}

	parcel := getParcel(w, r)
	parcel.Encode(http.StatusOK, list)
}

func (app *application) apiDeleteApp(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")

	// todo: should be adapter to delete app
	path, _ := os.Getwd()
	appPath := path + "/" + config.Config.Metadata.Path + "/" + appName
	err := os.RemoveAll(appPath)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func (app *application) apiPostApp(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")

	// Parse the JSON request body
	var requestData struct {
		Name string `json:"name"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// todo: should be adapter to rename app
	path, _ := os.Getwd()
	appPath := path + "/" + config.Config.Metadata.Path
	if requestData.Name != "" && appName != requestData.Name {
		err := os.Rename(appPath+"/"+appName, appPath+"/"+requestData.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) apiDuplicateApp(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")

	// todo: should be adapter to rename app
	path, _ := os.Getwd()
	appPath := path + "/" + config.Config.Metadata.Path

	err := Dir(appPath+"/"+appName, appPath+"/"+appName+"_copy")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) apiPostNewApp(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var requestData struct {
		Name string `json:"name"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// todo: should be an adapter to create app
	path, _ := os.Getwd()
	appPath := path + "/" + config.Config.Metadata.Path + "/" + requestData.Name

	subfolders := []string{config.Config.Metadata.Metasubpath, config.Config.Deployer.Entitysubpath}

	// Create the app folder
	err = os.Mkdir(appPath, 0755)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create the metadata and entity subfolders.
	for _, subfolder := range subfolders {
		path := appPath + "/" + subfolder
		err := os.MkdirAll(path, 0755)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

// func (app *application) postAppChatGPT(w http.ResponseWriter, r *http.Request) {

// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, "Error reading request body", http.StatusInternalServerError)
// 		return
// 	}

// 	var data ChatGPT
// 	if err := json.Unmarshal(body, &data); err != nil {
// 		http.Error(w, "Error unmarshalling request body", http.StatusBadRequest)
// 		return
// 	}

// 	// response, err := chatgpt.RenameApp(data.ChatGPT)
// 	// if err != nil {
// 	// 	if response != nil {
// 	// 		// Marshal map to JSON
// 	// 		jsonString, _ := json.Marshal(response)
// 	// 		if jsonString != nil {
// 	// 			http.Error(w, err.Error()+string(jsonString), http.StatusBadRequest)
// 	// 			return
// 	// 		}
// 	// 	}

// 	// 	http.Error(w, err.Error(), http.StatusBadRequest)
// 	// 	return
// 	// }

// 	http.Redirect(w, r, "/v1/", http.StatusSeeOther)
// }

func (app *application) apiRunApp(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")
	global.SetAppName(appName)

	catalogs, err := metaapi.ReadMetadataList(metadata.Catalog)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	events, err := metaapi.ReadMetadataList(metadata.Event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := struct {
		Catalogs []string
		Events   []string
		AppName  string
	}{
		catalogs,
		events,
		appName,
	}
	parcel := getParcel(w, r)
	parcel.Encode(http.StatusOK, data)
}
