package metadata

import (
	"time"
)

const IntegerType = "integer"
const NumberType = "number"
const StringType = "string"
const DatetimeType = "dateTime"
const BooleanType = "boolean"
const ReferenceType = "reference"
const TableType = "table"

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
	EntityName     string     `json:"entityName"`
	SearchCriteria []string   `json:"search"`
	Attributes     Attributes `json:"attributes"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

type Attributes map[string]Attribute

type Attribute map[string]any
