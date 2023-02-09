package adapter

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ukubenet/metadata-repository/models"
)

const Ext string = ".json"

type (
	// JSONCodec is a JSON implementation
	// of parcel.Decoder and parcel.Encoder
	JSONCodec struct {
		indent string
	}
)

// JSON returns a new JSON Encoder/Decoder
func JSON() *JSONCodec {
	return &JSONCodec{
		indent: "",
	}
}

// JSONIndent returns a new JSON Encoder/Decoder
// with the marshalled JSON indented by amt
func JSONIndent(amt int) *JSONCodec {
	return &JSONCodec{
		indent: strings.Repeat(" ", amt),
	}
}

// Read entity metadate from a json file
func (*JSONCodec) Read(entityName string, candidate *models.Entity) (err error) {

	file, err := os.Open(entityName + Ext)
	if err != nil {
		return
	}

	// close fi on exit and check for its returned error
	defer func() {
		if err := file.Close(); err != nil {
			return
		}
	}()

	jsonParser := json.NewDecoder(file)
	jsonParser.Decode(candidate)

	return
}

// Save the entity metadata to JSON file
func (jc *JSONCodec) Put(candidate *models.Entity) (err error) {
	var output []byte

	if jc.indent != "" {
		output, err = json.MarshalIndent(candidate, "", jc.indent)
	} else {
		output, err = json.Marshal(candidate)
	}

	if err != nil {
		return
	}

	os.WriteFile(candidate.EntityName+Ext, output, 0644)

	return
}

// Delete a json file containing entity metadata
func (*JSONCodec) Delete(entityName string) (err error) {
	err = os.Remove(entityName + Ext)
	return
}

// List json files containing entity metadata
func (*JSONCodec) List(list *[]string) (err error) {
	root, err := os.Getwd()
	filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// @todo pass const Ext below instead of ".json" string
		if filepath.Ext(info.Name()) == ".json" {
			*list = append(*list, path[len(root)+1:len(path)-5])
		}
		return nil
	})

	return
}
