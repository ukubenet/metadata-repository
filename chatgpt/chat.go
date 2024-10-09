package chatgpt

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/maiqingqiang/typechat-go"
	"github.com/ukubenet/metadata-repository/entity"
	entityapi "github.com/ukubenet/metadata-repository/entity/api"
	"github.com/ukubenet/metadata-repository/metadata"
	metaapi "github.com/ukubenet/metadata-repository/metadata/api"
)

func SaveNewEntity(entityType metadata.EntityType, name string, request string) (map[string]any, error) {
	meta, err := metaapi.ReadMetadata(entityType, name)
	if err != nil {
		return nil, err
	}

	scheme, err := json.Marshal(getRequestAttributesFromMeta(entityType, meta))
	if err != nil {
		return nil, err
	}

	response, err := chatGPTResponse(string(scheme), request)
	if err != nil {
		return nil, err
	}

	attributeValues, err := getEntityAttributeValuesFromResponse(*response, meta.GetStructedAttributes())
	if err != nil {
		return *response, err
	}

	e := new(entity.Entity)
	e.EntityName = name
	e.Attributes = attributeValues
	e.Identifier = uuid.New().String()

	err = entityapi.PutEntity(entityType, e)
	if err != nil {
		return *response, err
	}

	return *response, nil
}

func chatGPTResponse(scheme string, request string) (*map[string]any, error) {
	model, err := typechat.NewLanguageModel()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	translator := typechat.NewJsonTranslator[map[string]any](model, scheme, "Response")

	response, err := translator.Translate(request)
	if err != nil {
		log.Fatalf("translator.Translate Error: %v\n", err)
		return nil, err
	}

	// responseMap, ok := response.(map[string]any)
	// if !ok {
	// 	return nil, fmt.Errorf("response is not a map")
	// }

	if nestedResponse, ok := (*response)["Response"]; ok {
		if nestedMap, ok := nestedResponse.(map[string]any); ok {
			return &nestedMap, nil
		}
	}

	return response, nil

}

func getRequestAttributesFromMeta(entityType metadata.EntityType, meta *metadata.EntityMetadata) []any {

	requestAttributes := getRequestAttributesFromMetaAttributes(meta.Attributes)
	// if entityType == metadata.Event {
	// 	requestAttributes = append(requestAttributes, "time")
	// }

	return requestAttributes
}

func getRequestAttributesFromMetaAttributes(attributes metadata.Attributes) []any {
	requestAttributes := []any{}
	for name, value := range attributes {
		if value["type"] == "table" {
			columns := metadata.MapToAttributes(value["columns"].(map[string]interface{}))
			requestAttributes = append(requestAttributes, map[string]any{name: getRequestAttributesFromMetaAttributes(columns)})
		} else if value["type"] == "reference" {
			requestAttributes = append(requestAttributes, map[string]any{name: value["view"]})
		} else {
			requestAttributes = append(requestAttributes, name)
		}
	}

	return requestAttributes
}

func getEntityAttributeValuesFromResponse(response map[string]any, structedAttributes metadata.StructedAttributes) (entity.AttributeValues, error) {
	attributeValues := make(entity.AttributeValues)
	for name, value := range structedAttributes {
		if value.Type == metadata.ReferenceType {
			var err error
			attributeValues[name], err = findReference(response, name, value.Specs.(metadata.ReferenceSpecs))
			if err != nil {
				return nil, err
			}
		} else if value.Type == metadata.TableType {
			table := response[name]
			tableRows := []any{}
			for _, row := range table.([]any) {
				rowMap, ok := row.(map[string]any)
				if ok {
					tableRow, err := getEntityAttributeValuesFromResponse(rowMap, value.Specs.(metadata.TableSpecs).Columns)
					if err != nil {
						return nil, err
					}
					tableRows = append(tableRows, tableRow)
				}
			}
			attributeValues[name] = tableRows
		} else {
			attributeValues[name] = response[name]
		}
	}

	return attributeValues, nil
}

func findReference(response map[string]any, name string, refSpecs metadata.ReferenceSpecs) (map[string]any, error) {
	searchCriteria, ok := response[name]
	if !ok {
		searchCriteria = response
	}

	searchCriteriaMap, ok := searchCriteria.(map[string]any)
	if !ok {
		if len(refSpecs.View) == 1 {
			searchCriteriaString, ok := searchCriteria.(string)
			if ok {
				searchCriteriaMap = map[string]any{refSpecs.View[0]: searchCriteriaString}
			}
		} else {
			return nil, fmt.Errorf("search criteria for attribute %q is not a map", name)
		}
	}

	return entityapi.FindReference(refSpecs, searchCriteriaMap)
}
