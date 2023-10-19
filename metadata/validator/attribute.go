package metavalidator

import (
	"fmt"

	"github.com/ukubenet/metadata-repository/metadata"
	metastorage "github.com/ukubenet/metadata-repository/metadata/storage"
)

func ValidateAttributes(attributes metadata.Attributes) (err error) {
	for name, attribute := range attributes {
		err := ValidateAttribute(name, attribute)
		if err != nil {
			return err
		}
	}

	return
}

func ValidateAttribute(name string, attribute metadata.Attribute) (err error) {
	attrType, ok := attribute["type"]
	if !ok {
		return fmt.Errorf("attribute %q doesn't have type", name)
	}
	if _, ok := attrType.(string); !ok {
		return fmt.Errorf("type of attribute %q is not a string", name)
	}

	switch attrType.(string) {
	case metadata.ReferenceType:
		return validateReference(name, attribute)
	case metadata.TableType:
		return validateTable(name, attribute)
	case metadata.StringType:
	case metadata.IntegerType:
	case metadata.NumberType:
	case metadata.DatetimeType:
	case metadata.BooleanType:
		// do nothing
	default:
		return fmt.Errorf("type %q of attribute %q is not defined", attrType, name)
	}

	return
}

func validateReference(name string, attribute metadata.Attribute) (err error) {
	reference, ok := attribute["reference"]
	if !ok {
		return fmt.Errorf("reference attribute %q has missed reference property", name)
	}
	if _, ok := reference.(string); !ok {
		return fmt.Errorf("reference of attribute %q is not a string", name)
	}

	refEntity, err := metastorage.ReadMetadata(reference.(string))
	if err != nil {
		return fmt.Errorf(
			"error to read reference in attribute %q: %q",
			attribute["name"],
			err,
		)
	}

	attrview, ok := attribute["view"]
	if !ok {
		return fmt.Errorf("view of reference attribute %q is not present", name)
	}

	viewSlice, ok := attrview.([]interface{})
	if !ok {
		return fmt.Errorf("view of reference attribute %q is not a slice", name)
	}

	for _, viewElem := range viewSlice {
		view, ok := viewElem.(string)
		if !ok {
			return fmt.Errorf("view element %q is not a string", viewElem)
		}
		_, ok = refEntity.Attributes[view]
		if !ok {
			return fmt.Errorf("%q from view of reference attribute %q don't belong to reference entity", view, name)
		}
	}

	return
}

func validateTable(name string, attribute metadata.Attribute) (err error) {
	columns, ok := attribute["columns"]
	if !ok {
		return fmt.Errorf("table attribute %q has missed columns property", name)
	}
	metaColumns, ok := columns.(map[string]any)
	if !ok {
		return fmt.Errorf("columns property of attribute %q is malformed", name)
	}

	for columnName, columnMeta := range metaColumns {
		columnMetaAttr, ok := columnMeta.(map[string]any)
		if !ok {
			return fmt.Errorf("columns property of column %q is malformed", columnName)
		}
		ValidateAttribute(columnName, columnMetaAttr)
	}

	return
}
