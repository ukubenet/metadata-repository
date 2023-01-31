// Package storage provides mechanisms to write structs to different data repositories
package storage

import (
	"github.com/ukubenet/metadata-repository/models"
	"github.com/ukubenet/metadata-repository/storage/adapter"
)

type (
	// Reader implementations should decode values from a storage repository to a candidate
	Reader interface {
		Read(string, *models.Entity) error
	}

	// Replacer implementation should encode values from a candidate to a storage repository.
	Replacer interface {
		Put(candidate *models.Entity) error
	}

	// Adapter is a simple reference structure that
	// enables reading and writing from/to storage repository
	Adapter struct {
		factory *Factory
	}

	// Factory stores the implementation details of available adapters
	Factory struct {
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

// Use is a convience function to register storage adapter
func (f *Factory) Use(i interface{}) {
	if reader, ok := i.(Reader); ok {
		f.UseReader(reader)
	}

	if replacer, ok := i.(Replacer); ok {
		f.UseReplacer(replacer)
	}
}

func (f *Factory) CreateAdapter() *Adapter {
	return &Adapter{factory: f}
}

// Adapter

// Adapter insertion
func (a *Adapter) Put(c *models.Entity) error {
	inserter := a.factory.replacer

	return inserter.Put(c)
}

// Adapter reader
func (p *Adapter) Read(name string, c *models.Entity) (err error) {
	reader := p.factory.reader
	if err = reader.Read(name, c); err != nil {
		return
	}

	return
}

func CreateFactory() *Factory {
	factory := NewFactory()
	factory.Use(adapter.JSON())

	return factory
}
