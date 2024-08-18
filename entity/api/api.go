package entityapi

import (
	"errors"
	"fmt"

	"github.com/ukubenet/metadata-repository/entity"
	entitysearch "github.com/ukubenet/metadata-repository/entity/search"
	indexItem "github.com/ukubenet/metadata-repository/entity/search/item"
	storage "github.com/ukubenet/metadata-repository/entity/storage"
	entityvalidator "github.com/ukubenet/metadata-repository/entity/validator"
	"github.com/ukubenet/metadata-repository/metadata"
	metaapi "github.com/ukubenet/metadata-repository/metadata/api"
)

func ReadEntity(entityType metadata.EntityType, name string, identifier string) (*entity.Entity, error) {
	return storage.ReadEntity(entityType, name, identifier)
}

func PutEntity(entityType metadata.EntityType, entity *entity.Entity) error {

	if entity.EntityName == "" {
		return errors.New("entity name not defined")
	}
	if entity.Identifier == "" {
		return errors.New("entity identifier not defined")
	}
	if len(entity.Attributes) == 0 {
		return errors.New("entity attributes not defined")
	}
	if err := entityvalidator.ValidateAttributeValues(entityType, entity.EntityName, entity.Attributes); err != nil {
		return err
	}

	factoryWriter := storage.CreateFactory()
	adapter := factoryWriter.CreateAdapter()
	err := adapter.Put(entityType, entity)

	return err
}

func DeleteEntity(entityType metadata.EntityType, name string, identifier string) error {
	storage := storage.CreateFactory()
	adapter := storage.CreateAdapter()
	err := adapter.Delete(entityType, name, identifier)

	return err
}

func ReadEntities(entityType metadata.EntityType, name string) ([]entity.Entity, error) {
	return storage.ReadEntities(entityType, name)
}

func SearchEntities(entityType metadata.EntityType, name string, indexName string, criteria any) ([]indexItem.Key, error) {
	index := entitysearch.Indexes[entityType][name][indexName]

	list, err := index.Searcher.Search(criteria)

	return list, err
}

func ReadEntityTypes(entityType metadata.EntityType) ([]string, error) {
	storage := storage.CreateFactory()
	adapter := storage.CreateAdapter()

	list, err := adapter.TypeList(entityType)

	return list, err
}

func ReadReference(refSpecs *metadata.ReferenceSpecs, reference string) (entity.ReferenceValue, error) {
	refEntity, err := ReadEntity(refSpecs.EntityType, refSpecs.Reference, reference)
	if err != nil {
		return nil, err
	}
	view := make([]any, 0)
	for _, viewFieldName := range refSpecs.View {
		view = append(view, refEntity.Attributes[viewFieldName])
	}

	return entity.ReferenceValue{refEntity.Identifier: view}, nil
}

func FindSearchIndex(meta metadata.EntityMetadata, createria map[string]any) (*string, error) {
	var indexName *string
	for indName, indexMeta := range meta.SearchCriteria {
		foundIndex := true
		for _, fieldName := range indexMeta.Attributes {
			if _, ok := createria[fieldName]; !ok {
				foundIndex = false
				break
			}
		}

		if foundIndex {
			indexName = &indName
		}
		break
	}

	if indexName == nil {
		return indexName, fmt.Errorf("no search index for defined criteria")
	}

	return indexName, nil
}

func FindReference(refSpecs metadata.ReferenceSpecs, createria map[string]any) (entity.ReferenceValue, error) {
	refmeta, err := metaapi.ReadMetadata(refSpecs.EntityType, refSpecs.Reference)
	if err != nil {
		return nil, err
	}

	indexName, err := FindSearchIndex(*refmeta, createria)
	if err != nil {
		return nil, err
	}

	// indexValueMap := make(indexItem.ValueMap)
	// for _, fieldName := range refmeta.SearchCriteria[*indexName].Attributes {
	// 	indexValueMap[fieldName] = searchCriteriaMap[fieldName]
	// }

	list, err := SearchEntities(refSpecs.EntityType, refSpecs.Reference, *indexName, createria)
	if err != nil {
		return nil, err
	}
	if len(list) != 1 {
		return nil, fmt.Errorf("search by reference should find only 1 record")
	}

	return ReadReference(&refSpecs, string(list[0]))
}

func RetrieveReferenceByEntityId(attribute metadata.StructedAttribute, entityId string) (map[string]any, error) {
	refSpecs, ok := attribute.Specs.(metadata.ReferenceSpecs)
	if !ok {
		return nil, fmt.Errorf("meta type should be reference")
	}

	entityType := refSpecs.EntityType

	entity, err := ReadEntity(entityType, refSpecs.Reference, entityId)
	if err != nil {
		return nil, err
	}
	view := make(map[string]interface{})
	for _, viewFieldName := range refSpecs.View {
		view[viewFieldName] = entity.Attributes[viewFieldName]
	}

	return map[string]any{
		"reference": entityId,
		"type":      "reference",
		"view":      view,
	}, nil
}

// func GenerateReferenceAttributeValueForTableColumn(meta *metadata.EntityMetadata, attribute string, columnName string, reference string) (map[string]any, error) {
// 	if meta.Attributes[attribute]["type"] != "table" {
// 		return nil, fmt.Errorf("attribute %q should be a table type", attribute)
// 	}

// 	columnMetadata, ok := meta.Attributes[attribute]["rows"].(map[string]interface{})
// 	if !ok {
// 		return nil, fmt.Errorf("attribute %q does not have column metadata", attribute)
// 	}

// 	columnType, ok := columnMetadata[columnName].(map[string]interface{})["type"].(string)
// 	if !ok || columnType != "reference" {
// 		return nil, fmt.Errorf("column %q of attribute %q should be a reference type", columnName, attribute)
// 	}

// 	columnMetadataMap := columnMetadata[columnName].(map[string]interface{})

// 	viewFields := columnMetadataMap["view"]
// 	refType := columnMetadataMap["referenceType"].(string)
// 	referenceType, ok := metadata.EntityTypeMap[strings.ToLower(refType)]
// 	if !ok {
// 		return nil, fmt.Errorf("reference: %q, incorrect reference type: %q", columnMetadataMap["reference"].(string), refType)
// 	}

// 	refEntity, err := ReadEntity(referenceType, columnMetadataMap["reference"].(string), reference)
// 	if err != nil {
// 		return nil, err
// 	}
// 	view := make(map[string]interface{})
// 	for _, viewFieldName := range viewFields.([]interface{}) {
// 		view[viewFieldName.(string)] = refEntity.Attributes[viewFieldName.(string)]
// 	}

// 	return map[string]any{
// 		"reference": reference,
// 		"type":      "reference",
// 		"view":      view,
// 	}, nil
// }

func ReadReferences(metaSpecs *metadata.ReferenceSpecs) map[string]map[string]any {
	entities, err := ReadEntities(metaSpecs.EntityType, metaSpecs.Reference)
	if err != nil {
		panic(err)
	}

	result := make(map[string]map[string]any)
	for _, entity := range entities {
		view := make(map[string]any)
		for _, viewName := range metaSpecs.View {
			viewValue, ok := entity.Attributes[viewName]
			if ok {
				view[viewName] = viewValue
			}
		}

		result[entity.Identifier] = view
	}

	return result
}
