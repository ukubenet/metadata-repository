package metaapi

import (
	"errors"

	"github.com/ukubenet/metadata-repository/metadata"
	metastorage "github.com/ukubenet/metadata-repository/metadata/storage"
)

func ReadMetadata(name string) (*metadata.EntityMetadata, error) {
	entity := new(metadata.EntityMetadata)
	dbReader := metastorage.CreateFactory()
	adapter := dbReader.CreateAdapter()
	err := adapter.Read(name, entity)

	return entity, err
}

func ReadMetadataList() ([]string, error) {
	storage := metastorage.CreateFactory()
	adapter := storage.CreateAdapter()
	list := []string{}
	err := adapter.List(&list)

	return list, err
}

func PutMetadata(metadata *metadata.EntityMetadata) error {

	if metadata.EntityName == "" {
		return errors.New("entity name not defined")
	}
	if len(metadata.Attributes) == 0 {
		return errors.New("etity attributes not defined")
	}

	factoryWriter := metastorage.CreateFactory()
	adapter := factoryWriter.CreateAdapter()
	err := adapter.Put(metadata)

	return err
}

func DeleteMetadata(name string) error {

	storage := metastorage.CreateFactory()
	adapter := storage.CreateAdapter()
	err := adapter.Delete(name)

	return err
}
