package adapter

import (
	"encoding/json"
	"os"
	"time"

	"github.com/ukubenet/metadata-repository/entity"
)

func (jc *JSONCodec) RegisterRead(eventName string, identifier string) (records []*entity.Register, err error) {
	file, err := os.Open(jc.path + "/event/" + eventName + "/" + identifier + Ext)
	if err != nil {
		return nil, err
	}

	// close fi on exit and check for its returned error
	defer func() {
		if err := file.Close(); err != nil {
			return
		}
	}()

	jsonParser := json.NewDecoder(file)

	e := new(entity.Entity)
	if err = jsonParser.Decode(e); err != nil {
		return nil, err
	}

	return e.Registers, nil
}

func (jc *JSONCodec) RegisterPut(eventName string, identifier string, records []*entity.Register) error {
	var output []byte
	var err error

	if jc.indent != "" {
		output, err = json.MarshalIndent(records, "", jc.indent)
	} else {
		output, err = json.Marshal(records)
	}

	if err != nil {
		return err
	}

	err = os.WriteFile(jc.path+"/event/"+eventName+"/"+identifier+Ext, output, 0644)

	return err
}

func (jc *JSONCodec) RegisterDelete(eventName string, identifier string) error {
	err := os.Remove(jc.path + "/event/" + eventName + "/" + identifier + Ext)
	return err
}

func ReadState(registerType string, registerName string, dimensions []entity.AttributeValues, timestamp time.Time) (any, error) {
	return nil, nil
}
