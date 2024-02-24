package entityvalidator

import (
	"fmt"
	"strings"

	"github.com/ukubenet/metadata-repository/entity"
	entitystorage "github.com/ukubenet/metadata-repository/entity/storage"
	"github.com/ukubenet/metadata-repository/metadata"
	metaapi "github.com/ukubenet/metadata-repository/metadata/api"
	metavalidator "github.com/ukubenet/metadata-repository/metadata/validator"
)

func ValidateAttributeValues(entity string, values entity.AttributeValues) (err error) {
	meta, err := metaapi.ReadCatalogMetadata(entity)
	if err != nil {
		return fmt.Errorf("error reading meta of entity %q", entity)
	}

	if err = validateValues(values, meta.Attributes); err != nil {
		return err
	}

	return
}

func validateValues(values entity.AttributeValues, metaAttributes map[string]metadata.Attribute) (err error) {
	for name, value := range values {
		err = validateAttributeValue(name, value, metaAttributes[name])
		if err != nil {
			return err
		}
	}

	return
}

func validateAttributeValue(name string, value interface{}, meta metadata.Attribute) (err error) {
	if meta == nil {
		return fmt.Errorf("no such attribute %q", name)
	}

	err = metavalidator.ValidateAttribute(name, meta)
	if err != nil {
		return err
	}

	metatype := meta["type"].(string)

	switch meta["type"].(string) {
	case metadata.ReferenceType:
		return validateReference(name, value, meta)
	case metadata.TableType:
		return validateTable(name, value, meta)
	case metadata.StringType:
	case metadata.IntegerType:
	case metadata.NumberType:
	case metadata.DatetimeType:
	case metadata.BooleanType:
		// do nothing
	default:
		return fmt.Errorf("type %q of attribute %q is not defined", metatype, name)
	}

	return
}

func validateReference(name string, value any, meta metadata.Attribute) (err error) {
	referenceMap, ok := value.(map[string]any)
	if !ok {
		return fmt.Errorf("reference attribute %q is malformed", name)
	}

	reference, ok := referenceMap["reference"].(string)
	if !ok {
		return fmt.Errorf("reference of attribute %q does not exist or not a string", name)
	}

	referenceTypeString, ok := referenceMap["referenceType"].(string)
	if !ok {
		return fmt.Errorf("reference of attribute %q does not exist or not a string", name)
	}

	referenceType, ok := metadata.EntityTypeMap[strings.ToLower(referenceTypeString)]
	if !ok {
		return fmt.Errorf("reference type %q of attribute %q does not exist", referenceTypeString, name)
	}

	view, ok := referenceMap["view"]
	if !ok {
		return fmt.Errorf("view of reference attribute %q is not present", name)
	}

	if _, ok := view.(map[string]any); !ok {
		return fmt.Errorf("view of reference attribute %q is not a map", name)
	}

	refEntity, err := entitystorage.ReadEntity(referenceType, meta["reference"].(string), reference)
	if err != nil {
		return fmt.Errorf("error to read reference %q in attribute %q", reference, name)
	}

	for key, elem := range view.(map[string]any) {
		refValue, ok := refEntity.GetAttributes()[key]
		if !ok {
			return fmt.Errorf("attribute %q is not present in reference entity %q. entity attribute: %q", key, reference, name)
		}

		// @todo add comparison of other complex types such as references, etc
		if refValue != elem {
			return fmt.Errorf(
				"value %q from view attribute %q of reference attribute %q don't belong to reference entity. Value in ref entity: %q",
				elem,
				key,
				name,
				refValue,
			)
		}
	}

	return
}

func validateTable(name string, value interface{}, meta metadata.Attribute) (err error) {
	tableMap, ok := value.(map[string]any)
	if !ok {
		return fmt.Errorf("table attribute %q is not a string map", name)
	}

	columnsRaw, ok := tableMap["columns"]
	if !ok {
		return fmt.Errorf("table attribute %q doesn't have columns property", name)
	}

	columns, ok := columnsRaw.(map[string]any)
	if !ok {
		return fmt.Errorf("columns property of table attribute %q isn't a string map", name)
	}

	metaColumnsRaw, ok := meta["columns"]
	if !ok {
		return fmt.Errorf("table attribute %q doesn't have columns meta property", name)
	}

	metaColumnsMap, ok := metaColumnsRaw.(map[string]any)
	if !ok {
		return fmt.Errorf("columns meta property of table attribute %q isn't a string", name)
	}

	for columnName, columnValue := range columns {
		metaColumnRaw, ok := metaColumnsMap[columnName]
		if !ok {
			return fmt.Errorf("meta property for column %q of table attribute %q does not exist", columnName, name)
		}
		metaColumn, ok := metaColumnRaw.(map[string]any)
		if !ok {
			return fmt.Errorf("meta property for column %q of table attribute %q isn't a string map", columnName, name)
		}

		err = validateAttributeValue(columnName, columnValue, metaColumn)
		if err != nil {
			return err
		}
	}

	return
}
