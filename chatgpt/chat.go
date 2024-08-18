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

func SaveNewEntity(entityType metadata.EntityType, name string, request string) error {
	meta, err := metaapi.ReadMetadata(entityType, name)
	if err != nil {
		return err
	}

	scheme, err := json.Marshal(getRequestAttributesFromMeta(entityType, meta))
	if err != nil {
		return err
	}

	response, err := chatGPTResponse(string(scheme), request)
	if err != nil {
		return err
	}

	attributeValues, err := getEntityAttributeValuesFromResponse(*response, meta)
	if err != nil {
		return err
	}

	e := new(entity.Entity)
	e.EntityName = name
	e.Attributes = attributeValues
	e.Identifier = uuid.New().String()

	err = entityapi.PutEntity(entityType, e)
	if err != nil {
		return err
	}

	return nil
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

	// responseMap, ok := response.(*map[string]any)
	// if !ok {
	// 	return nil, fmt.Errorf("response is not a map")
	// }

	return response, nil
}

func getRequestAttributesFromMeta(entityType metadata.EntityType, meta *metadata.EntityMetadata) []any {

	requestAttributes := getRequestAttributesFromMetaAttributes(meta.Attributes)
	if entityType == metadata.Event {
		requestAttributes = append(requestAttributes, "time")
	}

	return requestAttributes
}

func getRequestAttributesFromMetaAttributes(attributes metadata.Attributes) []any {
	requestAttributes := []any{}
	for name, value := range attributes {
		if value["type"] == "table" {
			requestAttributes = append(requestAttributes, getRequestAttributesFromMetaAttributes(value["columns"].(metadata.Attributes)))
		} else if value["type"] == "reference" {
			requestAttributes = append(requestAttributes, map[string]any{name: value["view"]})
		} else {
			requestAttributes = append(requestAttributes, name)
		}
	}

	return requestAttributes
}

func getEntityAttributeValuesFromResponse(response map[string]any, meta *metadata.EntityMetadata) (entity.AttributeValues, error) {

	attributeValues := make(entity.AttributeValues)
	for name, value := range meta.GetStructedAttributes() {
		if value.Type == metadata.ReferenceType {
			var err error
			attributeValues[name], err = findReference(response, name, value.Specs.(metadata.ReferenceSpecs))
			if err != nil {
				return nil, err
			}
		} else {
			attributeValues[name] = response[name]
		}
	}

	return attributeValues, nil
}

func findReference(response map[string]any, name string, refSpecs metadata.ReferenceSpecs) (entity.ReferenceValue, error) {
	searchCriteria, ok := response[name]
	if !ok {
		return nil, fmt.Errorf("attribute %q is not in response", name)
	}

	searchCriteriaMap, ok := searchCriteria.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("search criteria for attribute %q is not a map", name)
	}

	return entityapi.FindReference(refSpecs, searchCriteriaMap)
}
