// Package metastorage provides mechanisms to write structs for entity metadata
package entitystorage

import (
	"github.com/ukubenet/metadata-repository/entity"
	"github.com/ukubenet/metadata-repository/metadata"
)

type (
	// Reader implementations should decode values from a storage repository to a candidate
	EntityReader interface {
		Read(entityType metadata.EntityType, entityName string, identifier string) (e *entity.Entity, err error)
	}

	// Replacer implementation should encode values from a candidate to a storage repository.
	EntityReplacer interface {
		Put(entityType metadata.EntityType, candidate *entity.Entity) error
	}

	// Delete implementations should delete entity
	EntityEraser interface {
		Delete(metadata.EntityType, string, string) error
	}

	// List implementations should show list of entities of specific type
	EntityLister interface {
		List(metadata.EntityType, string) ([]entity.Entity, error)
	}

	// List implementations should show list of entity types
	EntityTypeLister interface {
		TypeList(metadata.EntityType) ([]string, error)
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
func (a *Adapter) Put(entityType metadata.EntityType, c *entity.Entity) error {
	inserter := a.factory.replacer

	return inserter.Put(entityType, c)
}

// Adapter reader
func (p *Adapter) Read(entityType metadata.EntityType, name string, identifier string) (e *entity.Entity, err error) {
	reader := p.factory.reader
	if e, err = reader.Read(entityType, name, identifier); err != nil {
		return
	}

	return
}

// Adapter delete
func (p *Adapter) Delete(entityType metadata.EntityType, name string, identifier string) (err error) {
	eraser := p.factory.eraser
	if err = eraser.Delete(entityType, name, identifier); err != nil {
		return
	}

	return
}

// Adapter list
func (p *Adapter) List(entityType metadata.EntityType, name string) (list []entity.Entity, err error) {
	lister := p.factory.lister
	if list, err = lister.List(entityType, name); err != nil {
		return
	}

	return
}

// Adapter type list
func (p *Adapter) TypeList(entityType metadata.EntityType) (list []string, err error) {
	lister := p.factory.typeLister
	if list, err = lister.TypeList(entityType); err != nil {
		return
	}

	return
}

func CreateFactory() *EntityFactory {
	factory := NewFactory()
	factory.Use(getAdapter())

	return factory
}

func ReadEntity(entityType metadata.EntityType, name string, identifier string) (entity *entity.Entity, err error) {
	dbReader := CreateFactory()
	adapter := dbReader.CreateAdapter()
	entity, err = adapter.Read(entityType, name, identifier)

	return entity, err
}

func ReadEntities(entityType metadata.EntityType, name string) ([]entity.Entity, error) {
	storage := CreateFactory()
	adapter := storage.CreateAdapter()

	list, err := adapter.List(entityType, name)

	return list, err
}

func ReadEntityTypeList(entityType metadata.EntityType) (list []string, err error) {
	storage := CreateFactory()
	adapter := storage.CreateAdapter()

	list, err = adapter.TypeList(entityType)

	return list, err
}
