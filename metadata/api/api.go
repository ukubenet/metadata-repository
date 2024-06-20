package metaapi

import (
	"errors"

	"github.com/ukubenet/metadata-repository/metadata"
	metastorage "github.com/ukubenet/metadata-repository/metadata/storage"
	metavalidator "github.com/ukubenet/metadata-repository/metadata/validator"
)

func ReadMetadata(entityType metadata.EntityType, name string) (*metadata.EntityMetadata, error) {
	meta, err := metastorage.ReadMetadata(entityType, name)
	if err != nil {
		return meta, err
	}

	return meta, err
}

func ReadMetadataList(entityType metadata.EntityType) ([]string, error) {
	return metastorage.ReadMetadataList(entityType)
}

func PutMetadata(entityType metadata.EntityType, entitymeta *metadata.EntityMetadata) error {

	if entitymeta.EntityName == "" {
		return errors.New("entity name not defined")
	}
	if len(entitymeta.Attributes) == 0 {
		return errors.New("entity attributes not defined")
	}

	if entityType == metadata.Event {
		if _, ok := entitymeta.Attributes["EventTime"]; !ok {
			return errors.New("event entity must have EventTime attribute")
		}
		if entitymeta.Attributes["EventTime"]["type"] != metadata.DatetimeType {
			return errors.New("type of EventTime attribute must be dateTime")
		}
	}

	if err := metavalidator.ValidateAttributes(entitymeta.Attributes); err != nil {
		return err
	}

	return metastorage.PutMetadata(entityType, entitymeta)
}

func DeleteMetadata(entityType metadata.EntityType, name string) error {
	return metastorage.DeleteMetadata(entityType, name)
}
