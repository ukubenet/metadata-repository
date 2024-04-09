package entity

import (
	"time"

	"github.com/ukubenet/metadata-repository/metadata"
)

type Entity interface {
	GetID() string
	GetType() metadata.EntityType
	GetName() string
	GetAttributes() AttributeValues
}

type Metadata struct {
	EntityName string          `json:"entityName"`
	Identifier string          `json:"identifier"`
	Attributes AttributeValues `json:"attributes"`
}

type EventEntity struct {
	Metadata
	EventTime time.Time `json:"time"`
}

type CatalogEntity struct {
	Metadata
}

type AttributeValues map[string]any


func (e EventEntity) GetID() string {
	return e.Identifier
}
func (e EventEntity) GetName() string {
	return e.EntityName
}
func (e EventEntity) GetType() metadata.EntityType {
	return metadata.Event
}
func (e EventEntity) GetAttributes() AttributeValues {
	return e.Attributes
}

func (e CatalogEntity) GetID() string {
	return e.Identifier
}
func (e CatalogEntity) GetName() string {
	return e.EntityName
}
func (e CatalogEntity) GetType() metadata.EntityType {
	return metadata.Catalog
}
func (e CatalogEntity) GetAttributes() AttributeValues {
	return e.Attributes
}

func GetEntityInstance(entityType metadata.EntityType) (e Entity) {
	if metadata.Catalog == entityType {
		catalogEntity := new(CatalogEntity)
		return catalogEntity
	} else if metadata.Event == entityType {
		eventEntity := new(EventEntity)
		return eventEntity
	}

	return nil
}
