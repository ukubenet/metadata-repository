package entityapi

import (
	"errors"

	"github.com/ukubenet/metadata-repository/entity"
	entitysearch "github.com/ukubenet/metadata-repository/entity/search"
	indexItem "github.com/ukubenet/metadata-repository/entity/search/item"
	storage "github.com/ukubenet/metadata-repository/entity/storage"
	entityvalidator "github.com/ukubenet/metadata-repository/entity/validator"
	"github.com/ukubenet/metadata-repository/metadata"
)

func ReadEntity(entityType metadata.EntityType, name string, identifier string) (entity.Entity, error) {
	return storage.ReadEntity(entityType, name, identifier)
}

func PutEntity(entity entity.Entity) error {

	if entity.GetName() == "" {
		return errors.New("entity name not defined")
	}
	if entity.GetID() == "" {
		return errors.New("entity identifier not defined")
	}
	if len(entity.GetAttributes()) == 0 {
		return errors.New("entity attributes not defined")
	}
	if err := entityvalidator.ValidateAttributeValues(entity.GetType(), entity.GetName(), entity.GetAttributes()); err != nil {
		return err
	}

	factoryWriter := storage.CreateFactory()
	adapter := factoryWriter.CreateAdapter()
	err := adapter.Put(entity)

	return err
}

func DeleteEntity(entityType metadata.EntityType, name string, identifier string) error {
	storage := storage.CreateFactory()
	adapter := storage.CreateAdapter()
	err := adapter.Delete(entityType, name, identifier)

	return err
}

func ReadEntities(entityType metadata.EntityType, name string) ([]entity.Entity, error) {
	return storage.ReadEntities(entityType, name)
}

func SearchEntities(entityType metadata.EntityType, name string, index string, criteria any) ([]indexItem.Key, error) {
	indexTree := entitysearch.Indexes[entityType][name][index]

	list, err := indexTree.Searcher.Search(criteria)

	return list, err
}

func ReadEntityTypes(entityType metadata.EntityType) ([]string, error) {
	storage := storage.CreateFactory()
	adapter := storage.CreateAdapter()

	list, err := adapter.TypeList(entityType)

	return list, err
}
