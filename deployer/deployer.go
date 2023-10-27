// Package deployer provides mechanisms to deploy infrastructure to maintain entities
package deployer

import (
	"errors"

	"github.com/ukubenet/metadata-repository/deployer/adapter"
	"github.com/ukubenet/metadata-repository/metadata"
)

type (
	Deployer interface {
		Deploy(*metadata.EntityMetadata) error
	}

	Adapter struct {
		factory *Factory
	}

	Factory struct {
		deployer Deployer
	}
)

func NewFactory() *Factory {
	f := new(Factory)
	return f
}

func (f *Factory) UseDeployer(deployer Deployer) {
	f.deployer = deployer
}

func (f *Factory) Use(i interface{}) {
	if deployer, ok := i.(Deployer); ok {
		f.UseDeployer(deployer)
	}
}

func (f *Factory) CreateAdapter() *Adapter {
	return &Adapter{factory: f}
}

// Adapter

func (p *Adapter) Deploy(entitymeta *metadata.EntityMetadata) (err error) {
	deployer := p.factory.deployer
	if err = deployer.Deploy(entitymeta); err != nil {
		return
	}

	return
}

func CreateFactory() *Factory {
	factory := NewFactory()
	factory.Use(adapter.Local(get_path()))

	return factory
}

func DeployMetadata(entitymeta *metadata.EntityMetadata) error {
	if entitymeta.EntityName == "" {
		return errors.New("entity name not defined")
	}

	entity := new(metadata.EntityMetadata)
	dbReader := CreateFactory()
	adapter := dbReader.CreateAdapter()
	err := adapter.Deploy(entity)

	return err
}
