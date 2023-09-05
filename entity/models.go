package entity

import "time"

// type Type struct {
// 	Name         string `json:"name"`
// 	DataType     string `json:"dataType"`
// 	NotNull      bool   `json:"notNull"`
// 	DefaultValue string `json:"defaultValue"`
// }

type Entity struct {
	EntityName string          `json:"entityName"`
	Identifier string          `json:"Identifier"`
	Attributes AttributeValues `json:"attributes"`
	UpdatedAt  time.Time       `json:"updatedAt"`
}

type AttributeValues map[string]any

type ReferenceValue struct {
	Reference string          `json:"reference"`
	View      AttributeValues `json:"view"`
}
