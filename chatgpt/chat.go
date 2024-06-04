package chatgpt

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/maiqingqiang/typechat-go"
	"github.com/ukubenet/metadata-repository/entity"
	entityapi "github.com/ukubenet/metadata-repository/entity/api"
	indexItem "github.com/ukubenet/metadata-repository/entity/search/item"
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

	entityapi.SaveEntity(entityType, name, attributeValues)

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

	requestAttributes := []any{}
	for name, value := range meta.Attributes {
		if value["type"] == "reference" {
			requestAttributes = append(requestAttributes, map[string]any{name: value["view"]})
		} else {
			requestAttributes = append(requestAttributes, name)
		}
	}
	if entityType == metadata.Event {
		requestAttributes = append(requestAttributes, "time")
	}

	return requestAttributes
}

func getEntityAttributeValuesFromResponse(response map[string]any, meta *metadata.EntityMetadata) (entity.AttributeValues, error) {

	attributeValues := make(entity.AttributeValues)
	for name, value := range meta.Attributes {
		if value["type"] == "reference" {
			refType := value["referenceType"].(string)
			referenceType, ok := metadata.EntityTypeMap[strings.ToLower(refType)]
			if !ok {
				return nil, fmt.Errorf("reference: %q, incorrect reference type: %q", value["reference"].(string), refType)
			}

			refmeta, err := metaapi.ReadMetadata(referenceType, value["reference"].(string))
			if err != nil {
				return nil, err
			}

			searchCriteria, ok := response[name]
			if !ok {
				return nil, fmt.Errorf("attribute %q is not in response", name)
			}

			searchCriteriaMap, ok := searchCriteria.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("search criteria for attribute %q is not a map", name)
			}

			var indexName string = ""
			for indName, indexMeta := range refmeta.SearchCriteria {
				foundIndex := true
				for _, fieldName := range indexMeta.Attributes {
					if _, ok := searchCriteriaMap[fieldName]; !ok {
						foundIndex = false
						break
					}
				}

				if foundIndex {
					indexName = indName
				}
				break
			}

			if indexName == "" {
				return nil, fmt.Errorf("no search index for defined criteria")
			}

			indexValueMap := make(indexItem.ValueMap)
			for _, fieldName := range refmeta.SearchCriteria[indexName].Attributes {
				indexValueMap[fieldName] = searchCriteriaMap[fieldName]
			}

			list, err := entityapi.SearchEntities(referenceType, meta.Attributes[name]["reference"].(string), indexName, indexValueMap)
			if err != nil {
				return nil, err
			}
			if len(list) != 1 {
				return nil, fmt.Errorf("search by reference should find only 1 record. Attribute name: %q", name)
			}
			attributeValues[name], err = entityapi.GenerateReferenceAttributeValue(meta, name, string(list[0]))
			if err != nil {
				return nil, err
			}
		} else {
			attributeValues[name] = response[name]
		}
	}

	return attributeValues, nil
}
