package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/ukubenet/metadata-repository/config"
	global "github.com/ukubenet/metadata-repository/global"
	"github.com/ukubenet/metadata-repository/metadata"
	metaapi "github.com/ukubenet/metadata-repository/metadata/api"
)

func fullTemplatePath(slice []string, prefix string, suffix string) []string {
	for i, element := range slice {
		slice[i] = prefix + element + suffix
	}
	return slice
}

func executeMetadataTemplate(w http.ResponseWriter, tmplName string, tmplList []string, data interface{}) {
	path, _ := os.Getwd()
	tmplFiles := fullTemplatePath(tmplList, path+"/"+config.Config.Metadata.Path+"/"+config.Config.Metadata.Tmplsubpath+"/", ".tmpl")
	executeTemplate(w, tmplName, tmplFiles, data)
}

func (app *application) editMetadata(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")
	global.SetAppName(appName)

	entityType := params.ByName("type")
	entType, ok := metadata.EntityTypeMap[strings.ToLower(entityType)]
	if !ok {
		http.Error(w, "incorrect entity type: "+entityType, http.StatusBadRequest)
		return
	}

	name := params.ByName("name")
	entity, err := metaapi.ReadMetadata(entType, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmplData := struct {
		AppName        string
		EntityType     string
		Entity         *metadata.EntityMetadata
		AttributeTypes []string
	}{appName, entityType, entity, metadata.AttributeTypeList}

	executeMetadataTemplate(w, "metadata_edit", []string{"metadata_edit", "metadata_edit_form"}, tmplData)
}

func (app *application) newMetadata(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")
	global.SetAppName(appName)

	entityType := params.ByName("type")

	tmplData := struct {
		AppName        string
		EntityType     string
		AttributeTypes []string
	}{appName, entityType, metadata.AttributeTypeList}

	executeMetadataTemplate(w, "metadata_new", []string{"metadata_new", "metadata_edit_form"}, tmplData)

}

func (app *application) appConfiguration(w http.ResponseWriter, r *http.Request) {
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

	tmplData := struct {
		AppName  string
		Catalogs []string
		Events   []string
	}{
		appName,
		catalogs,
		events,
	}
	executeAppTemplate(w, "config", tmplData)
}

func (app *application) postMetadataChatGPT(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var data ChatGPT
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Error unmarshalling request body", http.StatusBadRequest)
		return
	}

	// response, err := chatgpt.RenameApp(data.ChatGPT)
	// if err != nil {
	// 	if response != nil {
	// 		// Marshal map to JSON
	// 		jsonString, _ := json.Marshal(response)
	// 		if jsonString != nil {
	// 			http.Error(w, err.Error()+string(jsonString), http.StatusBadRequest)
	// 			return
	// 		}
	// 	}

	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	http.Redirect(w, r, "/v1/", http.StatusSeeOther)
}
