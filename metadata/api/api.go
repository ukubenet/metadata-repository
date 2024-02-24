package metaapi

import (
	"errors"
	"fmt"

	"github.com/ukubenet/metadata-repository/metadata"
	metastorage "github.com/ukubenet/metadata-repository/metadata/storage"
	metavalidator "github.com/ukubenet/metadata-repository/metadata/validator"
)

func ReadCatalogMetadata(name string) (*metadata.EntityMetadata, error) {
	meta, err := metastorage.ReadMetadata(name, metadata.Catalog)
	if err != nil {
		return meta, err
	}

	return meta, err
}

func ReadCatalogMetadataList() ([]string, error) {
	return metastorage.ReadMetadataList(metadata.Catalog)
}

func PutCatalogMetadata(entitymeta *metadata.EntityMetadata) error {

	if entitymeta.EntityName == "" {
		return errors.New("entity name not defined")
	}
	if len(entitymeta.Attributes) == 0 {
		return errors.New("etity attributes not defined")
	}

	if err := metavalidator.ValidateAttributes(entitymeta.Attributes); err != nil {
		return err
	}

	return metastorage.PutMetadata(entitymeta, metadata.Catalog)
}

func DeleteCatalogMetadata(name string) error {
	return metastorage.DeleteMetadata(name, metadata.Catalog)
}

func ReadEventMetadata(name string) (*metadata.EntityMetadata, error) {
	meta, err := metastorage.ReadMetadata(name, metadata.Event)
	if err != nil {
		return meta, err
	}

	return meta, err
}

func ReadEventMetadataList() ([]string, error) {
	return metastorage.ReadMetadataList(metadata.Event)
}

func PutEventMetadata(entitymeta *metadata.EntityMetadata) error {

	if entitymeta.EntityName == "" {
		return errors.New("entity name not defined")
	}
	if len(entitymeta.Attributes) == 0 {
		return errors.New("etity attributes not defined")
	}

	if err := metavalidator.ValidateAttributes(entitymeta.Attributes); err != nil {
		return err
	}

	return metastorage.PutMetadata(entitymeta, metadata.Event)
}

func DeleteEventMetadata(name string) error {
	return metastorage.DeleteMetadata(name, metadata.Event)
}

func ReadMetadata(entityType metadata.EntityType, name string) (*metadata.EntityMetadata, error) {
	if entityType == metadata.Catalog {
		return ReadCatalogMetadata(name)
	} else if entityType == metadata.Event {
		return ReadEventMetadata(name)
	}

	return nil, fmt.Errorf("no such entity type %q", entityType)
}
