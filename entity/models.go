package entity

import "time"

type CatalogEntity struct {
	EntityName string          `json:"entityName"`
	Identifier string          `json:"identifier"`
	Attributes AttributeValues `json:"attributes"`
}

type EventEntity struct {
	Entity    CatalogEntity
	EventTime time.Time `json:"time"`
}

type Entity interface {
    CatalogEntity | EventEntity
}

type AttributeValues map[string]any
