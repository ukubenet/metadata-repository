// Package metastorage provides mechanisms to write structs for entity metadata
package metastorage

import (
	"errors"

	"github.com/ukubenet/metadata-repository/metadata"
)

type (
	// Reister Reader implementations should decode values from a storage repository to a candidate
	RegisterReader interface {
		RegisterRead(string, *metadata.RegisterMetadata) error
	}

	// Register Replacer implementation should encode values from a candidate to a storage repository.
	RegisterReplacer interface {
		RegisterPut(candidate *metadata.RegisterMetadata) error
	}

	// Register Delete implementations should delete entity
	RegisterEraser interface {
		RegisterDelete(string) error
	}

	// Register List implementations should show list of entities
	RegisterLister interface {
		RegisterList(*[]string) error
	}

	// Register Adapter is a simple reference structure that
	// enables reading and writing from/to storage repository
	RegisterAdapter struct {
		factory *RegisterFactory
	}

	// Factory stores the implementation details of available adapters
	RegisterFactory struct {
		lister   RegisterLister
		eraser   RegisterEraser
		replacer RegisterReplacer
		reader   RegisterReader
	}
)

// New Register Factory creates an instance of storage repository
func NewRegisterFactory() *RegisterFactory {
	f := new(RegisterFactory)
	return f
}

// UseRegisterReplacer registers register insertion or replacement implementation with the adapter factory
func (f *RegisterFactory) UseRegisterReplacer(replacer RegisterReplacer) {
	f.replacer = replacer
}

// UseRegisterReader registers register reader with the adapter factory
func (f *RegisterFactory) UseRegisterReader(reader RegisterReader) {
	f.reader = reader
}

// UseRegisterEraser registers register eraser with the adapter factory
func (f *RegisterFactory) UseRegisterEraser(eraser RegisterEraser) {
	f.eraser = eraser
}

// UseRegisterLister registers register lister with the adapter factory
func (f *RegisterFactory) UseRegisterLister(lister RegisterLister) {
	f.lister = lister
}

// Use is a convience function to register register storage adapter
func (f *RegisterFactory) UseRegister(i interface{}) {
	if reader, ok := i.(RegisterReader); ok {
		f.UseRegisterReader(reader)
	}

	if replacer, ok := i.(RegisterReplacer); ok {
		f.UseRegisterReplacer(replacer)
	}

	if eraser, ok := i.(RegisterEraser); ok {
		f.UseRegisterEraser(eraser)
	}

	if lister, ok := i.(RegisterLister); ok {
		f.UseRegisterLister(lister)
	}
}

func (f *RegisterFactory) CreateRegisterAdapter() *RegisterAdapter {
	return &RegisterAdapter{factory: f}
}

// Register Adapter

// Register Adapter replace
func (a *RegisterAdapter) RegisterPut(c *metadata.RegisterMetadata) error {
	inserter := a.factory.replacer

	return inserter.RegisterPut(c)
}

// Register Adapter reader
func (p *RegisterAdapter) RegisterRead(name string, c *metadata.RegisterMetadata) (err error) {
	reader := p.factory.reader
	if err = reader.RegisterRead(name, c); err != nil {
		return
	}

	return
}

// Register Adapter delete
func (p *RegisterAdapter) RegisterDelete(name string) (err error) {
	eraser := p.factory.eraser
	if err = eraser.RegisterDelete(name); err != nil {
		return
	}

	return
}

// Adapter list
func (p *RegisterAdapter) RegisterList(list *[]string) (err error) {
	lister := p.factory.lister
	if err = lister.RegisterList(list); err != nil {
		return
	}

	return
}

func CreateRegisterFactory(registerType metadata.RegisterType) *RegisterFactory {
	factory := NewRegisterFactory()
	factory.UseRegister(getRegisterAdapter(registerType))

	return factory
}

func ReadRegisterMetadata(registerType metadata.RegisterType, name string) (*metadata.RegisterMetadata, error) {
	entity := new(metadata.RegisterMetadata)
	dbReader := CreateRegisterFactory(registerType)
	adapter := dbReader.CreateRegisterAdapter()
	err := adapter.RegisterRead(name, entity)

	return entity, err
}

func ReadRegisterMetadataList(e metadata.RegisterType) ([]string, error) {
	storage := CreateRegisterFactory(e)
	adapter := storage.CreateRegisterAdapter()
	list := []string{}
	err := adapter.RegisterList(&list)

	return list, err
}

func PutRegisterMetadata(registerType metadata.RegisterType, registermeta *metadata.RegisterMetadata) error {

	if registermeta.RegisterName == "" {
		return errors.New("register name not defined")
	}
	if len(registermeta.Dimensions) == 0 {
		return errors.New("register dimensions not defined")
	}
	if len(registermeta.Fact) == 0 {
		return errors.New("register fact not defined")
	}
	// if registermeta.Source.Reference == "" {
	// 	return errors.New("register source not defined")
	// }


	factoryWriter := CreateRegisterFactory(registerType)
	adapter := factoryWriter.CreateRegisterAdapter()
	err := adapter.RegisterPut(registermeta)

	return err
}

func DeleteRegisterMetadata(registerType metadata.RegisterType, name string) error {

	storage := CreateRegisterFactory(registerType)
	adapter := storage.CreateRegisterAdapter()
	err := adapter.RegisterDelete(name)

	return err
}
