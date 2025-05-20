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
	Fact         any                   `json:"fact"`
	Source       ReferenceValue        `json:"source"`
	Next         *Register
	Previous     *Register
}

type RegisterState struct {
	RegisterName string
	RegisterType metadata.RegisterType
	Timestamp    time.Time
	Dimensions   AttributeValues
	InitialState any
	Head         *Register
	FinalState   any
	Tail         *Register
}
