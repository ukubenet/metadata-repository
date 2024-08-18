package metadata

const IntegerType = "integer"
const NumberType = "number"
const StringType = "string"
const DatetimeType = "dateTime"
const BooleanType = "boolean"
const ReferenceType = "reference"
const TableType = "table"

type EntityType int

const (
	Catalog EntityType = 1
	Event   EntityType = 2
)

func (e EntityType) String() string {
	return [...]string{"Catalog", "Event"}[e-1]
}

var (
	EntityTypeMap = map[string]EntityType{
		"catalog": Catalog,
		"event":   Event,
	}
)

var AttributeTypeList = []string{
	IntegerType,
	NumberType,
	StringType,
	DatetimeType,
	BooleanType,
	ReferenceType,
	TableType,
}

type EntityMetadata struct {
	EntityName     string           `json:"entityName"`
	SearchCriteria map[string]Index `json:"search"`
	Attributes     Attributes       `json:"attributes"`
}

type Attributes map[string]Attribute

type Attribute map[string]any

type StructedAttributes map[string]StructedAttribute

type StructedAttribute struct {
	Type  string `json:"type"`
	Specs interface{}
}

type ReferenceSpecs struct {
	View       []string   `json:"view"`
	EntityType EntityType `json:"entityType"`
	Reference  string     `json:"reference"`
}

type TableSpecs struct {
	Columns StructedAttributes `json:"columns"`
}

type PrimitiveType struct {
}

func GetStructedAttributes(attributes Attributes) StructedAttributes {
	structedAttributes := make(StructedAttributes)

	for name, attribute := range attributes {
		attrType, ok := attribute["type"]
		if !ok {
			panic("attribute doesn't have type")
		}
		switch attrType.(string) {
		case ReferenceType:
			structedAttributes[name] = StructedAttribute{
				Type: ReferenceType,
				Specs: ReferenceSpecs{
					Reference:  attribute["reference"].(string),
					EntityType: EntityTypeMap[attribute["referenceType"].(string)],
					View:       interfaceArayToStringArray(attribute["view"].([]interface{})),
				},
			}
		case TableType:
			structedAttributes[name] = StructedAttribute{
				Type: TableType,
				Specs: TableSpecs{
					Columns: GetStructedAttributes(mapToAttributes(attribute["columns"].(map[string]interface{}))),
				},
			}
		default:
			structedAttributes[name] = StructedAttribute{
				Type:  attrType.(string),
				Specs: PrimitiveType{},
			}
		}
	}

	return structedAttributes
}

func (entity *EntityMetadata) GetStructedAttributes() StructedAttributes {
	return GetStructedAttributes(entity.Attributes)
}

func interfaceArayToStringArray(interfaceArray []interface{}) []string {
	stringArray := make([]string, len(interfaceArray))
	for i, v := range interfaceArray {
		stringArray[i] = v.(string)
	}
	return stringArray
}

func mapToAttributes(m map[string]any) Attributes {
	attributes := make(Attributes)
	for key, value := range m {
		attributes[key] = mapToAttribute(value.(map[string]any))
	}
	return attributes
}

func mapToAttribute(m map[string]any) Attribute {
	attribute := make(Attribute)
	for key, value := range m {
		attribute[key] = value
	}
	return attribute
}
