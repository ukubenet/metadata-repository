package entityapi

import (
	"errors"

	"github.com/ukubenet/metadata-repository/entity"
	storage "github.com/ukubenet/metadata-repository/entity/storage"
	entityvalidator "github.com/ukubenet/metadata-repository/entity/validator"
)

func ReadEntity(name string, identifier string) (*entity.Entity, error) {
	entity, err := storage.ReadEntity(name, identifier)
	if err != nil {
		return entity, err
	}

	return entity, err
}

func PutEntity(entity *entity.Entity) error {

	if entity.EntityName == "" {
		return errors.New("entity name not defined")
	}
	if entity.Identifier == "" {
		return errors.New("entity identifier not defined")
	}
	if len(entity.Attributes) == 0 {
		return errors.New("entity attributes not defined")
	}
	if err := entityvalidator.ValidateAttributeValues(entity.EntityName, entity.Attributes); err != nil {
		return err
	}

	factoryWriter := storage.CreateFactory()
	adapter := factoryWriter.CreateAdapter()
	err := adapter.Put(entity)

	return err
}

func DeleteEntity(name string, identifier string) error {
	storage := storage.CreateFactory()
	adapter := storage.CreateAdapter()
	err := adapter.Delete(name, identifier)

	return err
}

func ReadEntities(name string) ([]entity.Entity, error) {
	storage := storage.CreateFactory()
	adapter := storage.CreateAdapter()
	list := []entity.Entity{}
	err := adapter.List(name, &list)

	return list, err
}

func ReadEntityTypes() ([]string, error) {
	storage := storage.CreateFactory()
	adapter := storage.CreateAdapter()
	list := []string{}
	err := adapter.TypeList(&list)

	return list, err
}
