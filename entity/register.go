package entity

import (
	"time"

	"github.com/ukubenet/metadata-repository/metadata"
)

type Register struct {
	RegisterName string                `json:"registerName"`
	RegisterType metadata.RegisterType `json:"registerType"`
	Timestamp    time.Time             `json:"timestamp"`
	Dimensions   AttributeValues       `json:"dimensions"`
	Facts        AttributeValues       `json:"facts"`
	Auxiliaries  AttributeValues       `json:"auxiliaries"`
	Source       ReferenceValue        `json:"source"`
}

type RegisterState struct {
	RegisterName string
	RegisterType metadata.RegisterType
	Timestamp    time.Time
	Dimensions   AttributeValues
	State        AttributeValues
	Auxiliaries  AttributeValues
}
