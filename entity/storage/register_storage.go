// Package metastorage provides mechanisms to write structs for entity metadata
package entitystorage

import (
	"time"

	"github.com/ukubenet/metadata-repository/entity"
)

type (
	EventRegisterReader interface {
		Read(eventName string, identifier string) (records []*entity.Register, err error)
	}

	EventRegisterReplacer interface {
		Put(eventName string, identifier string, records []*entity.Register) error
	}

	EventRegisterEraser interface {
		Delete(eventName string, identifier string) error
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
func (f *RegisterFactory) UseReplacer(replacer EventRegisterReplacer) {
	f.replacer = replacer
}

// UseReader registers reader with the adapter factory
func (f *RegisterFactory) UseReader(reader EventRegisterReader) {
	f.reader = reader
}

// UseEraser registers eraser with the adapter factory
func (f *RegisterFactory) UseEraser(eraser EventRegisterEraser) {
	f.eraser = eraser
}
// UseStateReader registers state reader with the adapter factory
func (f *RegisterFactory) UseStateReader(stateReader RegisterStateReader) {
	f.stateReader = stateReader
}

// Use is a convience function to register storage adapter
func (f *RegisterFactory) Use(i interface{}) {
	if reader, ok := i.(EventRegisterReader); ok {
		f.UseReader(reader)
	}

	if replacer, ok := i.(EventRegisterReplacer); ok {
		f.UseReplacer(replacer)
	}

	if eraser, ok := i.(EventRegisterEraser); ok {
		f.UseEraser(eraser)
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
func (a *RegisterAdapter) Put(eventName string, identifier string, records []*entity.Register) error {
	inserter := a.factory.replacer

	return inserter.Put(eventName, identifier, records)
}

// Adapter reader
func (p *RegisterAdapter) Read(eventName string, identifier string) (records []*entity.Register, err error) {
	reader := p.factory.reader
	if records, err = reader.Read(eventName, identifier); err != nil {
		return
	}

	return
}

// Adapter delete
func (p *RegisterAdapter) Delete(eventName string, identifier string) (err error) {
	eraser := p.factory.eraser
	if err = eraser.Delete(eventName, identifier); err != nil {
		return
	}

	return
}

// Adapter State Reader
func (p *RegisterAdapter) ReadState(registerType string, registerName string, dimensions []entity.AttributeValues, timestamp time.Time) (state any, err error) {
	stateReader := p.factory.stateReader
	if state, err = stateReader.ReadState(registerType, registerName,  dimensions, timestamp); err != nil {
		return
	}

	return
}


func CreateRegisterFactory() *RegisterFactory {
	factory := NewRegisterFactory()
	factory.Use(getAdapter())

	return factory
}

func ReadEventRegisters(eventName string, identifier string) (records []*entity.Register, err error) {
	dbReader := CreateRegisterFactory()
	adapter := dbReader.CreateRegisterAdapter()
	records, err = adapter.Read(eventName, identifier)

	return records, err
}


