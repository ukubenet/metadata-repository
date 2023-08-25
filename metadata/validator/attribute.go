package metavalidator

import (
	"fmt"
	"reflect"

	"github.com/ukubenet/metadata-repository/metadata"
	metastorage "github.com/ukubenet/metadata-repository/metadata/storage"
	"golang.org/x/exp/slices"
)

func ValidateAttributes(attributes []metadata.Attribute) (err error) {
	for i, attribute := range attributes {
		name, ok := attribute["name"]
		if !ok {
			return fmt.Errorf("attribute %d doesn't have name", i)
		}
		if _, ok := name.(string); !ok {
			return fmt.Errorf("name of attribute %d is not a string", i)
		}

		validateAttribute(attribute)
	}

	return
}

func validateAttribute(attribute metadata.Attribute) (err error) {
	name := attribute["name"].(string)
	attrType, ok := attribute["type"]
	if !ok {
		return fmt.Errorf("attribute %q doesn't have type", name)
	}
	if _, ok := attrType.(string); !ok {
		return fmt.Errorf("type of attribute %q is not a string", name)
	}

	switch attrType.(string) {
	case metadata.ReferenceType:
		return validateReference(attribute)
	case metadata.TableType:
		return validateTable(attribute)
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

func validateReference(attribute metadata.Attribute) (err error) {
	reference, ok := attribute["reference"]
	if !ok {
		return fmt.Errorf("reference attribute %q has missed reference property", attribute["name"])
	}
	if _, ok := reference.(string); !ok {
		return fmt.Errorf("reference of attribute %q is not a string", attribute["name"])
	}

	refEntity, err := metastorage.ReadMetadata(reference.(string))
	if err != nil {
		return fmt.Errorf(
			"error to read reference in attribute %q: %q",
			attribute["name"],
			err,
		)
	}

	view, ok := attribute["view"]
	if ok {
		if reflect.TypeOf(view).Kind() != reflect.Slice {
			return fmt.Errorf("view of reference attribute %q is not a slice", attribute["name"])
		}
		attributeNames := readAttributeNames(refEntity.Attributes)
		if !subslice(attribute["view"].([]string), attributeNames) {
			return fmt.Errorf(
				"some attributes from view of reference attribute %q don't belong to reference entity. Possible attributes: %q",
				attribute["name"],
				attributeNames,
			)
		}
	}

	return
}

func validateTable(attribute metadata.Attribute) (err error) {
	columns, ok := attribute["columns"]
	if !ok {
		return fmt.Errorf("table attribute %q has missed columns property", attribute["name"])
	}
	if _, ok := columns.([]metadata.Attribute); !ok {
		return fmt.Errorf("columns property of attribute %q is not a list of attributes", attribute["name"])
	}

	ValidateAttributes(columns.([]metadata.Attribute))

	return
}

func readAttributeNames(attributes []metadata.Attribute) []string {
	var attributeNames []string
	for _, attribute := range attributes {
		attributeNames = append(attributeNames, attribute["name"].(string))
	}

	return attributeNames
}

func subslice(s1 []string, s2 []string) bool {
	if len(s1) > len(s2) {
		return false
	}
	for _, e := range s1 {
		if !slices.Contains(s2, e) {
			return false
		}
	}
	return true
}
