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
	EntityName     string      `json:"entityName"`
	SearchCriteria []Attribute `json:"-"`
	Attributes     []Attribute `json:"attributes"`
	CreatedAt      time.Time   `json:"createdAt"`
	UpdatedAt      time.Time   `json:"updatedAt"`
}

type Attribute struct {
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Type struct {
	Name         string `json:"name"`
	DataType     string `json:"dataType"`
	NotNull      bool   `json:"notNull"`
	DefaultValue string `json:"defaultValue"`
}
