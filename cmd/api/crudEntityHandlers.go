package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/ukubenet/metadata-repository/chatgpt"
	"github.com/ukubenet/metadata-repository/config"
	"github.com/ukubenet/metadata-repository/entity"
	entityapi "github.com/ukubenet/metadata-repository/entity/api"
	global "github.com/ukubenet/metadata-repository/global"
	"github.com/ukubenet/metadata-repository/metadata"
	metaapi "github.com/ukubenet/metadata-repository/metadata/api"
)

var fns = template.FuncMap{
	"last": func(x int, a interface{}) bool {
		return x == reflect.ValueOf(a).Len()-1
	},
}

func executeTemplate(w http.ResponseWriter, tmplName string, tmplFile string, data interface{}) {
	tmpl, err := template.New(tmplName + ".tmpl").Funcs(fns).ParseFiles(tmplFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func executeEntityTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	path, _ := os.Getwd()
	tmplFile := path + "/" + config.Config.Metadata.Path + "/" + global.AppName + "/" + config.Config.Metadata.Tmplsubpath + "/" + tmplName + ".tmpl"
	executeTemplate(w, tmplName, tmplFile, data)
}

func createReferencesMap(meta *metadata.EntityMetadata) map[string]map[string]map[string]any {
	return populateReferencesMap(meta.GetStructedAttributes())
}

func populateReferencesMap(attributes metadata.StructedAttributes) map[string]map[string]map[string]any {
	references := make(map[string]map[string]map[string]any)

	for name, attribute := range attributes {
		if attribute.Type == metadata.ReferenceType {
			metaRefSpecs, ok := attribute.Specs.(metadata.ReferenceSpecs)
			if !ok {
				panic("metadata reference attribute error conversion")
			}

			references[name] = entityapi.ReadReferences(&metaRefSpecs)
		} else if attribute.Type == metadata.TableType {
			metaTableSpecs, ok := attribute.Specs.(metadata.TableSpecs)
			if !ok {
				panic("metadata table attribute error conversion")
			}

			tableReferences := populateReferencesMap(metaTableSpecs.Columns)
			for columnName, columnReferences := range tableReferences {
				references[name+"."+columnName] = columnReferences
			}
		}
	}

	return references
}

func getAttributesFromForm(r *http.Request, meta *metadata.EntityMetadata) (entity.AttributeValues, error) {
	metaAttributes := meta.GetStructedAttributes()
	AttributesValues := make(entity.AttributeValues)
	for key := range r.Form {
		var tableName string
		if idx := strings.IndexByte(key, '.'); idx >= 0 {
			tableName = key[:idx]
		}

		if tableName != "" && metaAttributes[tableName].Type == metadata.TableType {
			tableValues, err := getTableValuesFromForm(r, tableName, metaAttributes[tableName])
			if err != nil {
				return nil, err
			}
			AttributesValues[tableName] = tableValues
		} else if metaAttributes[key].Type == metadata.ReferenceType {
			refValue, err := entityapi.RetrieveReferenceByEntityId(metaAttributes[key], r.FormValue(key))
			if err != nil {
				return nil, err
			}
			AttributesValues[key] = refValue
		} else {
			AttributesValues[key] = r.FormValue(key)
		}
	}

	return AttributesValues, nil
}

func getTableValuesFromForm(r *http.Request, tableName string, meta metadata.StructedAttribute) ([]entity.AttributeValues, error) {
	tableSpecs, ok := meta.Specs.(metadata.TableSpecs)
	if !ok {
		return nil, fmt.Errorf("meta type should be table")
	}

	tableRows := make([]entity.AttributeValues, 0)

	for columnName, columnMeta := range tableSpecs.Columns {
		for rowIndex, formValue := range r.Form[tableName+"."+columnName] {
			if len(tableRows) <= rowIndex {
				tableRows = append(tableRows, make(entity.AttributeValues))
			}

			if columnMeta.Type == metadata.TableType {
				tableValues, err := getTableValuesFromForm(r, tableName+"."+columnName, columnMeta)
				if err != nil {
					return nil, err
				}
				tableRows[rowIndex][columnName] = tableValues
			} else if columnMeta.Type == metadata.ReferenceType {
				refValue, err := entityapi.RetrieveReferenceByEntityId(columnMeta, formValue)
				if err != nil {
					return nil, err
				}
				tableRows[rowIndex][columnName] = refValue
			} else {
				tableRows[rowIndex][columnName] = formValue
			}
		}
	}

	return tableRows, nil
}

func (app *application) editEntity(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")
	global.SetAppName(appName)

	name := params.ByName("name")
	identifier := params.ByName("identifier")
	entityType := params.ByName("type")
	entType, ok := metadata.EntityTypeMap[strings.ToLower(entityType)]
	if !ok {
		http.Error(w, "incorrect entity type: "+entityType, http.StatusBadRequest)
		return
	}

	e, err := entityapi.ReadEntity(entType, name, identifier)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	meta, err := metaapi.ReadMetadata(entType, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	references := createReferencesMap(meta)

	tmplData := struct {
		EntityType string
		Meta       *metadata.EntityMetadata
		Entity     entity.Entity
		References map[string]map[string]map[string]any
	}{entityType, meta, *e, references}

	executeEntityTemplate(w, "edit", tmplData)
}

func (app *application) newEntity(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")
	global.SetAppName(appName)

	name := params.ByName("name")
	entityType := params.ByName("type")
	entType, ok := metadata.EntityTypeMap[strings.ToLower(entityType)]
	if !ok {
		http.Error(w, "incorrect entity type: "+entityType, http.StatusBadRequest)
		return
	}

	meta, err := metaapi.ReadMetadata(entType, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	references := createReferencesMap(meta)

	tmplData := struct {
		EntityType string
		Meta       *metadata.EntityMetadata
		References map[string]map[string]map[string]any
	}{
		entityType,
		meta,
		references,
	}

	executeEntityTemplate(w, "new", tmplData)
}

func (app *application) listEntities(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")
	global.SetAppName(appName)

	name := params.ByName("name")
	entityType := params.ByName("type")
	entType, ok := metadata.EntityTypeMap[strings.ToLower(entityType)]
	if !ok {
		http.Error(w, "incorrect entity type: "+entityType, http.StatusBadRequest)
		return
	}

	meta, err := metaapi.ReadMetadata(entType, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	list, err := entityapi.ReadEntities(entType, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmplData := struct {
		EntityType string
		EntityName string
		Meta       *metadata.EntityMetadata
		List       []entity.Entity
	}{
		entityType,
		name,
		meta,
		list,
	}
	executeEntityTemplate(w, "list", tmplData)
}

func (app *application) postEntity(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")
	global.SetAppName(appName)

	name := params.ByName("name")
	identifier := params.ByName("identifier")
	entityType := params.ByName("type")
	entType, ok := metadata.EntityTypeMap[strings.ToLower(entityType)]
	if !ok {
		http.Error(w, "incorrect entity type: "+entityType, http.StatusBadRequest)
		return
	}

	r.ParseForm()
	if identifier == "" {
		identifier = r.FormValue("Identifier")
		r.Form.Del("Identifier")
	}

	meta, err := metaapi.ReadMetadata(entType, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	entity := new(entity.Entity)

	entity.EntityName = name
	entity.Identifier = identifier
	attributes, err := getAttributesFromForm(r, meta)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	entity.Attributes = attributes

	err = entityapi.PutEntity(entType, entity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/v1/list-view/"+entityType+"/"+name, http.StatusSeeOther)
}

type ChatGPT struct {
	ChatGPT string `json:"ChatGPT"`
}

func (app *application) postEntityChatGPT(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")
	global.SetAppName(appName)

	name := params.ByName("name")
	entityType := params.ByName("type")
	entType, ok := metadata.EntityTypeMap[strings.ToLower(entityType)]
	if !ok {
		http.Error(w, "incorrect entity type: "+entityType, http.StatusBadRequest)
		return
	}

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

	response, err := chatgpt.SaveNewEntity(entType, name, data.ChatGPT)
	if err != nil {
		if response != nil {
			// Marshal map to JSON
			jsonString, _ := json.Marshal(response)
			if jsonString != nil {
				http.Error(w, err.Error()+string(jsonString), http.StatusBadRequest)
				return
			}
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/v1/list-view/"+entityType+"/"+name, http.StatusSeeOther)
}
