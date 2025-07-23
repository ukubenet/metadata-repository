package metaapi

import (
	"errors"

	"github.com/ukubenet/metadata-repository/metadata"
	metastorage "github.com/ukubenet/metadata-repository/metadata/storage"
	metavalidator "github.com/ukubenet/metadata-repository/metadata/validator"
)

func ReadRegisterMetadata(registerType metadata.RegisterType, name string) (*metadata.RegisterMetadata, error) {
	meta, err := metastorage.ReadRegisterMetadata(registerType, name)
	if err != nil {
		return meta, err
	}

	return meta, err
}

func ReadRegisterMetadataList(registerType metadata.RegisterType) ([]string, error) {
	return metastorage.ReadRegisterMetadataList(registerType)
}

func PutRegisterMetadata(registerType metadata.RegisterType, registermeta *metadata.RegisterMetadata) error {

	if registermeta.RegisterName == "" {
		return errors.New("register name not defined")
	}
	if len(registermeta.Dimensions) == 0 {
		return errors.New("register dimensions not defined")
	}
	if registermeta.Facts == nil {
		return errors.New("register fact not defined")
	}
	// if registermeta.Source == nil {
	// 	return errors.New("register source not defined")
	// }

	if err := metavalidator.ValidateAttributes(registermeta.Dimensions); err != nil {
		return err
	}

	return metastorage.PutRegisterMetadata(registerType, registermeta)

}

func DeleteRegisterMetadata(registerType metadata.RegisterType, name string) error {
	return metastorage.DeleteRegisterMetadata(registerType, name)
}
