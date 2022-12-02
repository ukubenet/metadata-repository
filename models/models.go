package models

import "time"

type Entity struct {
	UUID           string            `json:"uuid"`
	EntityName     string            `json:"entityName"`
	CreatedAt      time.Time         `json:"createdAt"`
	UpdatedAt      time.Time         `json:"updatedAt"`
	SearchCriteria []EntityAttribute `json:"-"`
}

type Attribute struct {
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type EntityAttribute struct {
	UUID        string    `json:"uuid"`
	EntityId    string    `json:"entityId"`
	AttributeId string    `json:"attributeId"`
	Attribute   Attribute `json:"attribute"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
