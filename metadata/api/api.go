package metaapi

import (
	"errors"

	"github.com/ukubenet/metadata-repository/metadata"
	metastorage "github.com/ukubenet/metadata-repository/metadata/storage"
	metavalidator "github.com/ukubenet/metadata-repository/metadata/validator"
)

func ReadMetadata(name string) (*metadata.EntityMetadata, error) {
	meta, err := metastorage.ReadMetadata(name)
	if err != nil {
		return meta, err
	}

	return meta, err
}

func ReadMetadataList() ([]string, error) {
	return metastorage.ReadMetadataList()
}

func PutMetadata(entitymeta *metadata.EntityMetadata) error {

	if entitymeta.EntityName == "" {
		return errors.New("entity name not defined")
	}
	if len(entitymeta.Attributes) == 0 {
		return errors.New("etity attributes not defined")
	}

	if err := metavalidator.ValidateAttributes(entitymeta.Attributes); err != nil {
		return err
	}

	return metastorage.PutMetadata(entitymeta)
}

func DeleteMetadata(name string) error {
	return metastorage.DeleteMetadata(name)
}
