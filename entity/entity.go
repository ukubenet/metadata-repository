package entity

import (
	"github.com/ukubenet/metadata-repository/metadata"
	metaapi "github.com/ukubenet/metadata-repository/metadata/api"
)

type Entity struct {
	EntityName   string              `json:"entityName"`
	EntityType   metadata.EntityType `json:"entityType"`
	Identifier   string              `json:"identifier"`
	Attributes   AttributeValues     `json:"attributes"`
	Transactions map[string]any      `json:"transactions"`
	Registers    []*Register         `json:"registers"`
}

type AttributeValues map[string]any

type StructedAttributeValues map[string]StructedAttributeValue

type StructedAttributeValue struct {
	Value interface{}
	Meta  metadata.StructedAttribute `json:"meta"`
}

type ReferenceValue map[string]any // reference view map

type TableSpecs struct {
	Rows []StructedAttributeValues `json:"rows"`
}

type PrimitiveType struct {
}

func GetStructedAttributeValues(values AttributeValues, attributesMeta metadata.StructedAttributes) StructedAttributeValues {
	structedAttributes := make(StructedAttributeValues)

	for name, attribute := range attributesMeta {
		value, ok := values[name]
		if !ok {
			panic("value of attribute " + name + " is not present")
		}

		switch attribute.Type {
		case metadata.ReferenceType:
			if valueMap, ok := value.(map[string]interface{}); ok {
				if refValue, ok := valueMap["reference"].(string); ok {
					view, _ := valueMap["view"].([]any)
					structedAttributes[name] = StructedAttributeValue{
						Value: ReferenceValue{refValue: view},
						Meta:  attribute,
					}
				}
			}
		case metadata.TableType:
			if valueMap, ok := value.(map[string]interface{}); ok {
				rows, ok := valueMap["rows"].([]interface{})
				if !ok {
					panic("value of rows is not a slice")
				}
				structedRows := make([]StructedAttributeValues, len(rows))
				for i, row := range rows {
					if rowMap, ok := row.(map[string]interface{}); ok {
						structedRows[i] = GetStructedAttributeValues(AttributeValues(rowMap), attribute.Specs.(metadata.TableSpecs).Columns)
					}
				}
				structedAttributes[name] = StructedAttributeValue{
					Value: TableSpecs{
						Rows: structedRows,
					},
					Meta: attribute,
				}
			}
		default:
			structedAttributes[name] = StructedAttributeValue{
				Value: value,
				Meta:  attribute,
			}
		}
	}

	return structedAttributes
}

func (entity *Entity) GetStructedAttributes() StructedAttributeValues {
	meta, err := metaapi.ReadMetadata(entity.EntityType, entity.EntityName)
	if err != nil {
		panic(err)
	}
	return GetStructedAttributeValues(entity.Attributes, meta.GetStructedAttributes())
}

// https://stackoverflow.com/questions/45055953/interface-method-with-multiple-return-types
