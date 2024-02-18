package entityapi

import (
	"errors"

	"github.com/ukubenet/metadata-repository/entity"
	storage "github.com/ukubenet/metadata-repository/entity/storage"
	entityvalidator "github.com/ukubenet/metadata-repository/entity/validator"
	"github.com/ukubenet/metadata-repository/metadata"
)

func ReadCatalogEntity(name string, identifier string) (*entity.CatalogEntity, error) {
	entity, err := storage.ReadCatalogEntity(name, identifier)
	if err != nil {
		return entity, err
	}

	return entity, err
}

func PutCatalogEntity(entity *entity.CatalogEntity) error {

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

	factoryWriter := storage.CreateFactory(metadata.Catalog)
	adapter := factoryWriter.CreateAdapter()
	err := adapter.Put(entity)

	return err
}

func DeleteCatalogEntity(name string, identifier string) error {
	storage := storage.CreateFactory(metadata.Catalog)
	adapter := storage.CreateAdapter()
	err := adapter.Delete(name, identifier)

	return err
}

func ReadCatalogEntities(name string) ([]entity.CatalogEntity, error) {
	storage := storage.CreateFactory(metadata.Catalog)
	adapter := storage.CreateAdapter()
	list := []entity.CatalogEntity{}
	err := adapter.List(name, &list)

	return list, err
}

func ReadCatalogEntityTypes() ([]string, error) {
	storage := storage.CreateFactory(metadata.Catalog)
	adapter := storage.CreateAdapter()
	list := []string{}
	err := adapter.TypeList(&list)

	return list, err
}
