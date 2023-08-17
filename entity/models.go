package entity

import "time"

// type Type struct {
// 	Name         string `json:"name"`
// 	DataType     string `json:"dataType"`
// 	NotNull      bool   `json:"notNull"`
// 	DefaultValue string `json:"defaultValue"`
// }

type Entity struct {
	EntityName string           `json:"entityName"`
	Identifier string           `json:"Identifier"`
	Attributes []AttributeValue `json:"attributes"`
	UpdatedAt  time.Time        `json:"updatedAt"`
}	

type AttributeValue map[string]interface{}
