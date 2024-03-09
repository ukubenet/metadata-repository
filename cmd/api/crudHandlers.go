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

func (app *application) editCatalogEntity(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	name := params.ByName("name")
	identifier := params.ByName("identifier")

	e, err := entityapi.ReadEntity(metadata.Catalog, name, identifier)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path, _ := os.Getwd()
	tmplFile := path + config.Config.Metadata.Tmplpath + "catalog_edit.tmpl"
	tmpl, err := template.New("catalog_edit.tmpl").ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}

	var references = make(map[string]map[string]map[string]string)
	meta, _ := metaapi.ReadMetadata(e.GetType(), e.GetName())
	for name, value := range e.GetAttributes() {
		metaAttribute := meta.Attributes[name]
		if metaAttribute["type"] != "reference" {
			continue
		}

		referenceMap, ok := value.(map[string]any)
		if !ok {
			http.Error(w, fmt.Sprintf("reference attribute %q is malformed", name), http.StatusBadRequest)
			return
		}

		// @todo why do we retrieve referenceType from entity. It should be metadata
		referenceTypeString, ok := referenceMap["referenceType"].(string)
		if !ok {
			http.Error(w, fmt.Sprintf("reference type of attribute %q does not exist or not a string", name), http.StatusBadRequest)
			return
		}

		referenceType, ok := metadata.EntityTypeMap[strings.ToLower(referenceTypeString)]
		if !ok {
			http.Error(w, fmt.Sprintf("reference type %q of attribute %q does not match (catalog or event)", referenceTypeString, name), http.StatusBadRequest)
			return
		}

		refEntities, err := entityapi.ReadEntities(referenceType, metaAttribute["reference"].(string))
		if err != nil {
			http.Error(w, fmt.Sprintf("error to read reference %q in attribute %q", referenceType, metaAttribute["reference"].(string)), http.StatusBadRequest)
		}

		var refEntityMap = make(map[string]map[string]string)
		for _, refEntity := range refEntities {
			refIdentifier := refEntity.GetID()
			var refViewMap = make(map[string]string)
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

	tmplData := struct {
		Entity     entity.Entity
		References map[string]map[string]map[string]string
	}{
		e,
		references,
	}

	err = tmpl.Execute(w, tmplData)
	if err != nil {
		panic(err)
	}
}

func (app *application) newCatalogEntity(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	name := params.ByName("name")

	meta, err := metaapi.ReadCatalogMetadata(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	path, _ := os.Getwd()
	tmplFile := path + config.Config.Metadata.Tmplpath + "catalog_new.tmpl"
	tmpl, err := template.New("catalog_new.tmpl").ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}

	var references = make(map[string]map[string]map[string]string)
	for name, value := range meta.Attributes {
		metaAttribute := meta.Attributes[name]
		if value["type"] != "reference" {
			continue
		}

		refEntities, err := entityapi.ReadEntities(metadata.Catalog, metaAttribute["reference"].(string)) // todo support for Event and Catalog
		if err != nil {
			http.Error(w, fmt.Sprintf("error to read reference %q in attribute %q", metadata.Catalog, metaAttribute["reference"].(string)), http.StatusBadRequest)
		}

		var refEntityMap = make(map[string]map[string]string)
		for _, refEntity := range refEntities {
			refIdentifier := refEntity.GetID()
			var refViewMap = make(map[string]string)
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

	tmplData := struct {
		Meta       *metadata.EntityMetadata
		References map[string]map[string]map[string]string
	}{
		meta,
		references,
	}

	err = tmpl.Execute(w, tmplData)
	if err != nil {
		panic(err)
	}
}

func (app *application) postCatalogEntity(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	name := params.ByName("name")
	identifier := params.ByName("identifier")
	r.ParseForm()
	if identifier == "" {
		identifier = r.FormValue("Identifier")
		r.Form.Del("Identifier")
	}

	e := new(entity.CatalogEntity)
	e.EntityName = name
	e.Identifier = identifier

	meta, err := metaapi.ReadCatalogMetadata(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	AttributesValues := make(entity.AttributeValues)
	for key := range r.Form {
		if meta.Attributes[key]["type"] == "reference" {
			viewFields := meta.Attributes[key]["view"]
			refEntity, err := entityapi.ReadEntity(metadata.Catalog, meta.Attributes[key]["reference"].(string), r.FormValue(key)) // @todo fix it reference should be a parameter either catalog or event
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			view := make(map[string]interface{})
			for _, viewFieldName := range viewFields.([]interface{}) {
				view[viewFieldName.(string)] = refEntity.GetAttributes()[viewFieldName.(string)]
			}

			AttributesValues[key] = map[string]any{
				"reference":     r.FormValue(key),
				"type":          "reference",
				"referenceType": "Catalog", // @todo fix it reference should be a parameter either catalog or event
				"view":          view,
			}
		} else {
			AttributesValues[key] = r.FormValue(key)
		}
	}
	e.Attributes = AttributesValues

	err = entityapi.PutEntity(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/v1/catalog-view/"+name, http.StatusSeeOther)
}

func (app *application) listCatalogEntities(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	name := params.ByName("name")

	meta, err := metaapi.ReadCatalogMetadata(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	list, err := entityapi.ReadEntities(metadata.Catalog, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path, _ := os.Getwd()
	tmplFile := path + config.Config.Metadata.Tmplpath + "catalog_list.tmpl"

	tmpl, err := template.New("catalog_list.tmpl").Funcs(fns).ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}

	tmplData := struct {
		EntityName string
		Meta       *metadata.EntityMetadata
		List       []entity.Entity
	}{
		name,
		meta,
		list,
	}
	err = tmpl.Execute(w, tmplData)
	if err != nil {
		panic(err)
	}
}

func (app *application) editEventEntity(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	name := params.ByName("name")
	identifier := params.ByName("identifier")

	entity, err := entityapi.ReadEntity(metadata.Event, name, identifier)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path, _ := os.Getwd()
	tmplFile := path + config.Config.Metadata.Tmplpath + "event_edit.tmpl"
	tmpl, err := template.New("event_edit.tmpl").ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, entity)
	if err != nil {
		panic(err)
	}
}

func (app *application) newEventEntity(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	name := params.ByName("name")

	meta, err := metaapi.ReadEventMetadata(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	path, _ := os.Getwd()
	tmplFile := path + config.Config.Metadata.Tmplpath + "event_new.tmpl"
	tmpl, err := template.New("event_new.tmpl").ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, meta)
	if err != nil {
		panic(err)
	}
}

func (app *application) postEventEntity(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	name := params.ByName("name")
	identifier := params.ByName("identifier")

	r.ParseForm()
	if identifier == "" {
		identifier = r.FormValue("Identifier")
		r.Form.Del("Identifier")
	}

	e := new(entity.EventEntity)
	e.EntityName = name
	e.Identifier = identifier

	value := r.FormValue("EventTime")
	EventTime, err := time.Parse(time.RFC3339, value+":00Z")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	e.EventTime = EventTime
	r.Form.Del("EventTime")

	AttributesValues := make(entity.AttributeValues)
	for key := range r.Form {
		AttributesValues[key] = r.FormValue(key)
	}
	e.Attributes = AttributesValues

	err = entityapi.PutEntity(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/v1/event-view/"+name, http.StatusSeeOther)
}

func (app *application) listEventEntities(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	name := params.ByName("name")

	meta, err := metaapi.ReadEventMetadata(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	list, err := entityapi.ReadEntities(metadata.Event, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path, _ := os.Getwd()
	tmplFile := path + config.Config.Metadata.Tmplpath + "event_list.tmpl"

	tmpl, err := template.New("event_list.tmpl").ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}

	tmplData := struct {
		EntityName string
		Meta       *metadata.EntityMetadata
		List       []entity.Entity
	}{
		name,
		meta,
		list,
	}

	err = tmpl.Execute(w, tmplData)
	if err != nil {
		panic(err)
	}
}
