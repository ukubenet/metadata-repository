package entityapi

import (
	"time"

	"github.com/ukubenet/metadata-repository/entity"
	storage "github.com/ukubenet/metadata-repository/entity/storage"
)

func CommitEvent(eventName string, eventIdentifier string, registers []*entity.Register) error {
	factoryWriter := storage.CreateRegisterFactory()
	adapter := factoryWriter.CreateRegisterAdapter()
	return adapter.Put(eventName, eventIdentifier, registers)
}

func RollbackEvent(eventName string, eventIdentifier string) error {

	return nil
}

func ReadEventRegisters(eventName string, eventIdentifier string) ([]*entity.Register, error) {
	return nil, nil
}
func ReadState(registerType string, registerName string, dimensions []entity.AttributeValues, timestamp time.Time) (any, error) {
	return nil, nil
}
