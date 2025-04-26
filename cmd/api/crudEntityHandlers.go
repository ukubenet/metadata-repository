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
	"add": func(a, b int) int { return a + b },
	"inSlice": func(slice []string, element string) bool {
		for _, s := range slice {
			if s == element {
				return true
			}
		}
		return false
	},
	"convertToStringSlice": func(rawValues []interface{}) []string {
		values := make([]string, len(rawValues))
		for i, v := range rawValues {
			values[i] = v.(string)
		}
		return values
	},
}

func executeTemplate(w http.ResponseWriter, tmplName string, tmplFiles []string, data interface{}) {
	tmpl, err := template.New(tmplName + ".tmpl").Funcs(fns).ParseFiles(tmplFiles...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func executeEntityTemplate(w http.ResponseWriter, tmplName string, tmplList []string, data interface{}) {
	path, _ := os.Getwd()
	tmplFiles := fullTemplatePath(tmplList, path+"/"+config.Config.Metadata.Path+"/"+global.AppName+"/"+config.Config.Metadata.Tmplsubpath+"/", ".tmpl")
	tmplFiles = append(tmplFiles, path+"/"+config.Config.Metadata.Path+"/"+config.Config.Metadata.Tmplsubpath+"/"+"header"+".tmpl")
	tmplFiles = append(tmplFiles, path+"/"+config.Config.Metadata.Path+"/"+config.Config.Metadata.Tmplsubpath+"/"+"run_header"+".go.tmpl")
	executeTemplate(w, tmplName, tmplFiles, data)
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
		if strings.HasPrefix(key, "transaction") {
			continue
		}
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

func getTransactionsFromForm(r *http.Request, meta *metadata.EntityMetadata) (map[string]any, error) {
	metaAttributes := meta.Transactions
	transactions := make(map[string]any)
	for key := range metaAttributes {
		value := make(map[string]any)
		oldReference := "transaction['" + key + "']['old_entity_reference']"
		newReference := "transaction['" + key + "']['new_entity_reference']"
		oldChange := "transaction['" + key + "']['old_change']"
		newChange := "transaction['" + key + "']['new_change']"
		oldReferenceValue := r.FormValue(oldReference)
		if oldReferenceValue != "" {
			value["old_entity_reference"] = oldReferenceValue
		}
		newReferenceValue := r.FormValue(newReference)
		if newReferenceValue != "" {
			value["new_entity_reference"] = newReferenceValue
		}
		oldChangeValue := r.FormValue(oldChange)
		if oldChangeValue != "" {
			value["old_change"] = oldChangeValue
		}
		newChangeValue := r.FormValue(newChange)
		if newChangeValue != "" {
			value["new_change"] = newChangeValue
		}

		transactions[key] = value
	}

	return transactions, nil
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

func (app *application) appRun(w http.ResponseWriter, r *http.Request) {
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
	// balances, err := metaapi.ReadMetadataList(metadata.Balance)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	tmplData := struct {
		AppName  string
		Catalogs []string
		Events   []string
		Balances []string
	}{
		appName,
		catalogs,
		events,
		nil, //balances,
	}
	executeAppTemplate(w, "run", tmplData)
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
		AppName    string
		EntityType string
		Meta       *metadata.EntityMetadata
		Entity     entity.Entity
		References map[string]map[string]map[string]any
	}{appName, entityType, meta, *e, references}

	if entType == metadata.Event {
		executeEntityTemplate(w, "edit", []string{"edit", "on_submit", strings.ToLower(name) + "_transaction"}, tmplData)
	} else {
		executeEntityTemplate(w, "edit", []string{"edit", "on_submit", "no_transaction"}, tmplData)
	}
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
		AppName    string
		EntityType string
		Meta       *metadata.EntityMetadata
		References map[string]map[string]map[string]any
	}{
		appName,
		entityType,
		meta,
		references,
	}

	if entType == metadata.Event {
		executeEntityTemplate(w, "new", []string{"new", "on_submit", strings.ToLower(name) + "_transaction"}, tmplData)
	} else {
		executeEntityTemplate(w, "new", []string{"new", "on_submit", "no_transaction"}, tmplData)
	}
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
		AppName,
		EntityType string
		EntityName string
		Meta       *metadata.EntityMetadata
		List       []entity.Entity
	}{
		appName,
		entityType,
		name,
		meta,
		list,
	}
	executeEntityTemplate(w, "list", []string{"list"}, tmplData)
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

	if entType == metadata.Event {
		transactions, err := getTransactionsFromForm(r, meta)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		entity.Transactions = transactions
	}

	err = entityapi.PutEntity(entType, entity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/v1/list-view/"+appName+"/"+entityType+"/"+name, http.StatusSeeOther)
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

	http.Redirect(w, r, "/v1/list-view/"+appName+"/"+entityType+"/"+name, http.StatusSeeOther)
}
