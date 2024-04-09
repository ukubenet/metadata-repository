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