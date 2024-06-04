package entity

type Entity struct {
	EntityName string          `json:"entityName"`
	Identifier string          `json:"identifier"`
	Attributes AttributeValues `json:"attributes"`
}

type AttributeValues map[string]any
