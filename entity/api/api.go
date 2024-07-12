package entityapi

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/ukubenet/metadata-repository/entity"
	entitysearch "github.com/ukubenet/metadata-repository/entity/search"
	indexItem "github.com/ukubenet/metadata-repository/entity/search/item"
	storage "github.com/ukubenet/metadata-repository/entity/storage"
	entityvalidator "github.com/ukubenet/metadata-repository/entity/validator"
	"github.com/ukubenet/metadata-repository/metadata"
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

func SearchEntities(entityType metadata.EntityType, name string, index string, criteria any) ([]indexItem.Key, error) {
	indexTree := entitysearch.Indexes[entityType][name][index]

	list, err := indexTree.Searcher.Search(criteria)

	return list, err
}

func ReadEntityTypes(entityType metadata.EntityType) ([]string, error) {
	storage := storage.CreateFactory()
	adapter := storage.CreateAdapter()

	list, err := adapter.TypeList(entityType)

	return list, err
}

func SaveEntity(entityType metadata.EntityType, name string, attributeValues entity.AttributeValues) error {
	e := new(entity.Entity)
	e.EntityName = name
	e.Attributes = attributeValues
	e.Identifier = uuid.New().String()

	return PutEntity(entityType, e)
}

func GenerateReferenceAttributeValue(meta *metadata.EntityMetadata, name string, reference string) (map[string]any, error) {
	if meta.Attributes[name]["type"] != "reference" {
		return nil, fmt.Errorf("attribute %q should be a reference type", name)
	}
	viewFields := meta.Attributes[name]["view"]
	refType := meta.Attributes[name]["referenceType"].(string)
	referenceType, ok := metadata.EntityTypeMap[strings.ToLower(refType)]
	if !ok {
		return nil, fmt.Errorf("reference: %q, incorrect reference type: %q", meta.Attributes[name]["reference"].(string), refType)
	}

	refEntity, err := ReadEntity(referenceType, meta.Attributes[name]["reference"].(string), reference)
	if err != nil {
		return nil, err
	}
	view := make(map[string]interface{})
	for _, viewFieldName := range viewFields.([]interface{}) {
		view[viewFieldName.(string)] = refEntity.Attributes[viewFieldName.(string)]
	}

	return map[string]any{
		"reference": reference,
		"type":      "reference",
		"view":      view,
	}, nil
}

func GenerateReferenceAttributeValueForTableColumn(meta *metadata.EntityMetadata, attribute string, columnName string, reference string) (map[string]any, error) {
	if meta.Attributes[attribute]["type"] != "table" {
		return nil, fmt.Errorf("attribute %q should be a table type", attribute)
	}

	columnMetadata, ok := meta.Attributes[attribute]["columns"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("attribute %q does not have column metadata", attribute)
	}

	columnType, ok := columnMetadata[columnName].(map[string]interface{})["type"].(string)
	if !ok || columnType != "reference" {
		return nil, fmt.Errorf("column %q of attribute %q should be a reference type", columnName, attribute)
	}

	columnMetadataMap := columnMetadata[columnName].(map[string]interface{})

	viewFields := columnMetadataMap["view"]
	refType := columnMetadataMap["referenceType"].(string)
	referenceType, ok := metadata.EntityTypeMap[strings.ToLower(refType)]
	if !ok {
		return nil, fmt.Errorf("reference: %q, incorrect reference type: %q", columnMetadataMap["reference"].(string), refType)
	}

	refEntity, err := ReadEntity(referenceType, columnMetadataMap["reference"].(string), reference)
	if err != nil {
		return nil, err
	}
	view := make(map[string]interface{})
	for _, viewFieldName := range viewFields.([]interface{}) {
		view[viewFieldName.(string)] = refEntity.Attributes[viewFieldName.(string)]
	}

	return map[string]any{
		"reference": reference,
		"type":      "reference",
		"view":      view,
	}, nil
}
