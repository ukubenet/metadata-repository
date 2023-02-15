package models

import (
	"time"
)

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
