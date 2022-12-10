package models

import (
	"database/sql"
	"time"
)

// Models is the wrapper for database
type Models struct {
	DB DBModel
}

// NewModels returns models with db pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

type Entity struct {
	UUID           string      `json:"uuid"`
	EntityName     string      `json:"entityName"`
	SearchCriteria []Attribute `json:"-"`
	Attributes     []Attribute `json:"attributes"`
	CreatedAt      time.Time   `json:"createdAt"`
	UpdatedAt      time.Time   `json:"updatedAt"`
}

type Attribute struct {
	UUID       string    `json:"uuid"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	EntityUuid string    `json:"entityUuid"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type Type struct {
	UUID         string `json:"uuid"`
	Name         string `json:"name"`
	DataType     string `json:"dataType"`
	NotNull      bool   `json:"notNull"`
	DefaultValue string `json:"defaultValue"`
}

type EntityAttribute struct {
	UUID        string    `json:"uuid"`
	EntityId    string    `json:"entityId"`
	AttributeId string    `json:"attributeId"`
	Attribute   Attribute `json:"attribute"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
