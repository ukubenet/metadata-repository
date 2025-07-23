// Package metastorage provides mechanisms to write structs for entity metadata
package entitystorage

import (
	"time"

	"github.com/ukubenet/metadata-repository/entity"
)

type (
	EventRegisterReader interface {
		RegisterRead(eventName string, identifier string) (records []*entity.Register, err error)
	}

	EventRegisterReplacer interface {
		RegisterPut(eventName string, identifier string, records *entity.Register) error
	}

	EventRegisterEraser interface {
		RegisterDelete(eventName string, identifier string) error
	}

	RegisterStateReader interface {
		ReadState(registerType string, registerName string, dimensions []entity.AttributeValues, timestamp time.Time) (any, error)
	}

	// Adapter is a simple reference structure that
	// enables reading and writing from/to storage repository
	RegisterAdapter struct {
		factory *RegisterFactory
	}

	// Factory stores the implementation details of available adapters
	RegisterFactory struct {
		replacer    EventRegisterReplacer
		reader      EventRegisterReader
		eraser      EventRegisterEraser
		stateReader RegisterStateReader
	}
)

// NewFactory creates an instance of storage repository
func NewRegisterFactory() *RegisterFactory {
	f := new(RegisterFactory)
	return f
}

// UseReplacer registers insertion or replacement implementation with the adapter factory
func (f *RegisterFactory) UseRegisterReplacer(replacer EventRegisterReplacer) {
	f.replacer = replacer
}

// UseReader registers reader with the adapter factory
func (f *RegisterFactory) UseRegisterReader(reader EventRegisterReader) {
	f.reader = reader
}

// UseEraser registers eraser with the adapter factory
func (f *RegisterFactory) UseRegisterEraser(eraser EventRegisterEraser) {
	f.eraser = eraser
}

// UseStateReader registers state reader with the adapter factory
func (f *RegisterFactory) UseStateReader(stateReader RegisterStateReader) {
	f.stateReader = stateReader
}

// Use is a convience function to register storage adapter
func (f *RegisterFactory) Use(i interface{}) {
	if reader, ok := i.(EventRegisterReader); ok {
		f.UseRegisterReader(reader)
	}

	if replacer, ok := i.(EventRegisterReplacer); ok {
		f.UseRegisterReplacer(replacer)
	}

	if eraser, ok := i.(EventRegisterEraser); ok {
		f.UseRegisterEraser(eraser)
	}

	if stateReader, ok := i.(RegisterStateReader); ok {
		f.UseStateReader(stateReader)
	}
}

func (f *RegisterFactory) CreateRegisterAdapter() *RegisterAdapter {
	return &RegisterAdapter{factory: f}
}

// Adapter

// Adapter replace
func (a *RegisterAdapter) RegisterPut(eventName string, identifier string, records *entity.Register) error {
	inserter := a.factory.replacer

	return inserter.RegisterPut(eventName, identifier, records)
}

// Adapter reader
func (p *RegisterAdapter) RegisterRead(eventName string, identifier string) (records []*entity.Register, err error) {
	reader := p.factory.reader
	if records, err = reader.RegisterRead(eventName, identifier); err != nil {
		return
	}

	return
}

// Adapter delete
func (p *RegisterAdapter) RegisterDelete(eventName string, identifier string) (err error) {
	eraser := p.factory.eraser
	if err = eraser.RegisterDelete(eventName, identifier); err != nil {
		return
	}

	return
}

// Adapter State Reader
func (p *RegisterAdapter) ReadState(registerType string, registerName string, dimensions []entity.AttributeValues, timestamp time.Time) (state any, err error) {
	stateReader := p.factory.stateReader
	if state, err = stateReader.ReadState(registerType, registerName, dimensions, timestamp); err != nil {
		return
	}

	return
}

func CreateRegisterFactory() *RegisterFactory {
	factory := NewRegisterFactory()
	factory.Use(getRegisterAdapter())

	return factory
}

func ReadEventRegisters(eventName string, identifier string) (records []*entity.Register, err error) {
	dbReader := CreateRegisterFactory()
	adapter := dbReader.CreateRegisterAdapter()
	records, err = adapter.RegisterRead(eventName, identifier)

	return records, err
}
