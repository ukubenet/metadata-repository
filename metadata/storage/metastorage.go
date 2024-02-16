// Package metastorage provides mechanisms to write structs for entity metadata
package metastorage

import (
	"errors"

	"github.com/ukubenet/metadata-repository/metadata"
)

type (
	// Reader implementations should decode values from a storage repository to a candidate
	Reader interface {
		Read(string, *metadata.EntityMetadata) error
	}

	// Replacer implementation should encode values from a candidate to a storage repository.
	Replacer interface {
		Put(candidate *metadata.EntityMetadata) error
	}

	// Delete implementations should delete entity
	Eraser interface {
		Delete(string) error
	}

	// List implementations should show list of entities
	Lister interface {
		List(*[]string) error
	}

	// Adapter is a simple reference structure that
	// enables reading and writing from/to storage repository
	Adapter struct {
		factory *Factory
	}

	// Factory stores the implementation details of available adapters
	Factory struct {
		lister   Lister
		eraser   Eraser
		replacer Replacer
		reader   Reader
	}
)

// NewFactory creates an instance of storage repository
func NewFactory() *Factory {
	f := new(Factory)
	return f
}

// UseReplacer registers insertion or replacement implementation with the adapter factory
func (f *Factory) UseReplacer(replacer Replacer) {
	f.replacer = replacer
}

// UseReader registers reader with the adapter factory
func (f *Factory) UseReader(reader Reader) {
	f.reader = reader
}

// UseEraser registers eraser with the adapter factory
func (f *Factory) UseEraser(eraser Eraser) {
	f.eraser = eraser
}

// UseLister registers lister with the adapter factory
func (f *Factory) UseLister(lister Lister) {
	f.lister = lister
}

// Use is a convience function to register storage adapter
func (f *Factory) Use(i interface{}) {
	if reader, ok := i.(Reader); ok {
		f.UseReader(reader)
	}

	if replacer, ok := i.(Replacer); ok {
		f.UseReplacer(replacer)
	}

	if eraser, ok := i.(Eraser); ok {
		f.UseEraser(eraser)
	}

	if lister, ok := i.(Lister); ok {
		f.UseLister(lister)
	}
}

func (f *Factory) CreateAdapter() *Adapter {
	return &Adapter{factory: f}
}

// Adapter

// Adapter replace
func (a *Adapter) Put(c *metadata.EntityMetadata) error {
	inserter := a.factory.replacer

	return inserter.Put(c)
}

// Adapter reader
func (p *Adapter) Read(name string, c *metadata.EntityMetadata) (err error) {
	reader := p.factory.reader
	if err = reader.Read(name, c); err != nil {
		return
	}

	return
}

// Adapter delete
func (p *Adapter) Delete(name string) (err error) {
	eraser := p.factory.eraser
	if err = eraser.Delete(name); err != nil {
		return
	}

	return
}

// Adapter list
func (p *Adapter) List(list *[]string) (err error) {
	lister := p.factory.lister
	if err = lister.List(list); err != nil {
		return
	}

	return
}

func CreateFactory(e metadata.EntityType) *Factory {
	factory := NewFactory()
	factory.Use(getAdapter(e))

	return factory
}

func ReadMetadata(name string, e metadata.EntityType) (*metadata.EntityMetadata, error) {
	entity := new(metadata.EntityMetadata)
	dbReader := CreateFactory(e)
	adapter := dbReader.CreateAdapter()
	err := adapter.Read(name, entity)

	return entity, err
}

func ReadMetadataList(e metadata.EntityType) ([]string, error) {
	storage := CreateFactory(e)
	adapter := storage.CreateAdapter()
	list := []string{}
	err := adapter.List(&list)

	return list, err
}

func PutMetadata(entitymeta *metadata.EntityMetadata, e metadata.EntityType) error {

	if entitymeta.EntityName == "" {
		return errors.New("entity name not defined")
	}
	if len(entitymeta.Attributes) == 0 {
		return errors.New("etity attributes not defined")
	}

	factoryWriter := CreateFactory(e)
	adapter := factoryWriter.CreateAdapter()
	err := adapter.Put(entitymeta)

	return err
}

func DeleteMetadata(name string, e metadata.EntityType) error {

	storage := CreateFactory(e)
	adapter := storage.CreateAdapter()
	err := adapter.Delete(name)

	return err
}
