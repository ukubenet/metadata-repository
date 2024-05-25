package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/ukubenet/metadata-repository/config"
	"github.com/ukubenet/metadata-repository/entity"
	entityapi "github.com/ukubenet/metadata-repository/entity/api"
	"github.com/ukubenet/metadata-repository/metadata"
	metaapi "github.com/ukubenet/metadata-repository/metadata/api"
)

var fns = template.FuncMap{
	"last": func(x int, a interface{}) bool {
		return x == reflect.ValueOf(a).Len()-1
	},
}

func executeTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	path, _ := os.Getwd()
	tmplFile := path + config.Config.Metadata.Tmplpath + tmplName + ".tmpl"
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

func createReferencesMap(meta *metadata.EntityMetadata) (map[string]map[string]map[string]string, error) {
	references := make(map[string]map[string]map[string]string)

	for name, value := range meta.Attributes {
		metaAttribute := meta.Attributes[name]
		if value["type"] != "reference" {
			continue
		}
		referenceTypeString, ok := metaAttribute["referenceType"].(string)
		if !ok {
			return nil, fmt.Errorf("reference type of attribute %q does not exist or not a string (not define in metadata)", name)
		}

		referenceType, ok := metadata.EntityTypeMap[strings.ToLower(referenceTypeString)]
		if !ok {
			return nil, fmt.Errorf("reference type %q of attribute %q does not exist (not define in metadata)", referenceTypeString, name)
		}

		refEntities, err := entityapi.ReadEntities(referenceType, metaAttribute["reference"].(string))
		if err != nil {
			return nil, err
		}

		refEntityMap := make(map[string]map[string]string)
		for _, refEntity := range refEntities {
			refIdentifier := refEntity.GetID()
			refViewMap := make(map[string]string)
			for _, refViewName := range metaAttribute["view"].([]interface{}) {
				refViewValue, ok := refEntity.GetAttributes()[refViewName.(string)].(string)
				if ok {
					refViewMap[refViewName.(string)] = refViewValue
				}
			}

			refEntityMap[refIdentifier] = refViewMap
		}
		references[name] = refEntityMap
	}

	return references, nil
}

func getAttributesFromForm(r *http.Request, meta *metadata.EntityMetadata) (entity.AttributeValues, error) {
	AttributesValues := make(entity.AttributeValues)
	for key := range r.Form {
		if meta.Attributes[key]["type"] == "reference" {
			viewFields := meta.Attributes[key]["view"]
			refType := meta.Attributes[key]["referenceType"].(string)
			referenceType, ok := metadata.EntityTypeMap[strings.ToLower(refType)]
			if !ok {
				return nil, fmt.Errorf("reference: %q, incorrect reference type: %q", meta.Attributes[key]["reference"].(string), refType)
			}

			refEntity, err := entityapi.ReadEntity(referenceType, meta.Attributes[key]["reference"].(string), r.FormValue(key))
			if err != nil {
				return nil, err
			}
			view := make(map[string]interface{})
			for _, viewFieldName := range viewFields.([]interface{}) {
				view[viewFieldName.(string)] = refEntity.GetAttributes()[viewFieldName.(string)]
			}

			AttributesValues[key] = map[string]any{
				"reference": r.FormValue(key),
				"type":      "reference",
				"view":      view,
			}
		} else {
			AttributesValues[key] = r.FormValue(key)
		}
	}

	return AttributesValues, nil
}

func (app *application) editEntity(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

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

	references, err := createReferencesMap(meta)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmplData := struct {
		EntityType string
		Entity     entity.Entity
		References map[string]map[string]map[string]string
	}{entityType, e, references}

	executeTemplate(w, "edit", tmplData)
}

func (app *application) newEntity(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
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

	references, err := createReferencesMap(meta)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmplData := struct {
		EntityType string
		Meta       *metadata.EntityMetadata
		References map[string]map[string]map[string]string
	}{
		entityType,
		meta,
		references,
	}

	executeTemplate(w, "new", tmplData)
}

func (app *application) listEntities(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
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
	executeTemplate(w, "list", tmplData)
}

func (app *application) postEntity(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

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

	if entType == metadata.Event {
		eventEntity := new(entity.EventEntity)

		value := r.FormValue("EventTime")
		EventTime, err := time.Parse(time.RFC3339, value+":00Z")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		eventEntity.EventTime = EventTime
		r.Form.Del("EventTime")

		eventEntity.EntityName = name
		eventEntity.Identifier = identifier
		attributes, err := getAttributesFromForm(r, meta)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		eventEntity.Attributes = attributes

		err = entityapi.PutEntity(eventEntity)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		entityCatalog := new(entity.CatalogEntity)

		entityCatalog.EntityName = name
		entityCatalog.Identifier = identifier
		attributes, err := getAttributesFromForm(r, meta)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		entityCatalog.Attributes = attributes

		err = entityapi.PutEntity(entityCatalog)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	http.Redirect(w, r, "/v1/list-view/"+entityType+"/"+name, http.StatusSeeOther)
}
