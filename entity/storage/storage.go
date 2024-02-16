// Package metastorage provides mechanisms to write structs for entity metadata
package entitystorage

import (
	"github.com/ukubenet/metadata-repository/entity"
)

type (
	// Reader implementations should decode values from a storage repository to a candidate
	EntityReader interface {
		Read(string, string, *entity.CatalogEntity) error
	}

	// Replacer implementation should encode values from a candidate to a storage repository.
	EntityReplacer interface {
		Put(candidate *entity.CatalogEntity) error
	}

	// Delete implementations should delete entity
	EntityEraser interface {
		Delete(string, string) error
	}

	// List implementations should show list of entities of specific type
	EntityLister interface {
		List(string, *[]entity.CatalogEntity) error
	}

	// List implementations should show list of entity types
	EntityTypeLister interface {
		TypeList(*[]string) error
	}

	// Adapter is a simple reference structure that
	// enables reading and writing from/to storage repository
	Adapter struct {
		factory *EntityFactory
	}

	// Factory stores the implementation details of available adapters
	EntityFactory struct {
		lister     EntityLister
		eraser     EntityEraser
		replacer   EntityReplacer
		reader     EntityReader
		typeLister EntityTypeLister
	}
)

// NewFactory creates an instance of storage repository
func NewFactory() *EntityFactory {
	f := new(EntityFactory)
	return f
}

// UseReplacer registers insertion or replacement implementation with the adapter factory
func (f *EntityFactory) UseReplacer(replacer EntityReplacer) {
	f.replacer = replacer
}

// UseReader registers reader with the adapter factory
func (f *EntityFactory) UseReader(reader EntityReader) {
	f.reader = reader
}

// UseEraser registers eraser with the adapter factory
func (f *EntityFactory) UseEraser(eraser EntityEraser) {
	f.eraser = eraser
}

// UseLister registers lister with the adapter factory
func (f *EntityFactory) UseLister(lister EntityLister) {
	f.lister = lister
}

// UseTypeLister registers lister with the adapter factory
func (f *EntityFactory) UseTypeLister(typeLister EntityTypeLister) {
	f.typeLister = typeLister
}

// Use is a convience function to register storage adapter
func (f *EntityFactory) Use(i interface{}) {
	if reader, ok := i.(EntityReader); ok {
		f.UseReader(reader)
	}

	if replacer, ok := i.(EntityReplacer); ok {
		f.UseReplacer(replacer)
	}

	if eraser, ok := i.(EntityEraser); ok {
		f.UseEraser(eraser)
	}

	if lister, ok := i.(EntityLister); ok {
		f.UseLister(lister)
	}

	if typeLister, ok := i.(EntityTypeLister); ok {
		f.UseTypeLister(typeLister)
	}
}

func (f *EntityFactory) CreateAdapter() *Adapter {
	return &Adapter{factory: f}
}

// Adapter

// Adapter replace
func (a *Adapter) Put(c *entity.CatalogEntity) error {
	inserter := a.factory.replacer

	return inserter.Put(c)
}

// Adapter reader
func (p *Adapter) Read(name string, identifier string, c *entity.CatalogEntity) (err error) {
	reader := p.factory.reader
	if err = reader.Read(name, identifier, c); err != nil {
		return
	}

	return
}

// Adapter delete
func (p *Adapter) Delete(name string, identifier string) (err error) {
	eraser := p.factory.eraser
	if err = eraser.Delete(name, identifier); err != nil {
		return
	}

	return
}

// Adapter list
func (p *Adapter) List(name string, list *[]entity.CatalogEntity) (err error) {
	lister := p.factory.lister
	if err = lister.List(name, list); err != nil {
		return
	}

	return
}

// Adapter type list
func (p *Adapter) TypeList(list *[]string) (err error) {
	lister := p.factory.typeLister
	if err = lister.TypeList(list); err != nil {
		return
	}

	return
}

func CreateFactory() *EntityFactory {
	factory := NewFactory()
	factory.Use(getAdapter())

	return factory
}

func ReadEntity(name string, identifier string) (*entity.CatalogEntity, error) {
	entity := new(entity.CatalogEntity)
	dbReader := CreateFactory()
	adapter := dbReader.CreateAdapter()
	err := adapter.Read(name, identifier, entity)

	return entity, err
}
