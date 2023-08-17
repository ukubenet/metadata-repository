package adapter

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ukubenet/metadata-repository/metadata"
)

const Ext string = ".json"

type (
	// JSONCodec is a JSON implementation
	// of parcel.Decoder and parcel.Encoder
	JSONCodec struct {
		indent string
		path   string
	}
)

// JSON returns a new JSON Encoder/Decoder
func JSON(path string) *JSONCodec {
	return &JSONCodec{
		indent: "",
		path:   path,
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
func (jc *JSONCodec) Read(entityName string, candidate *metadata.EntityMetadata) (err error) {

	file, err := os.Open(jc.path + entityName + Ext)
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
func (jc *JSONCodec) Put(candidate *metadata.EntityMetadata) (err error) {
	var output []byte

	if jc.indent != "" {
		output, err = json.MarshalIndent(candidate, "", jc.indent)
	} else {
		output, err = json.Marshal(candidate)
	}

	if err != nil {
		return
	}

	os.WriteFile(jc.path+candidate.EntityName+Ext, output, 0644)

	return
}

// Delete a json file containing entity metadata
func (jc *JSONCodec) Delete(entityName string) (err error) {
	err = os.Remove(jc.path + entityName + Ext)
	return
}

// List json files containing entity metadata
func (jc *JSONCodec) List(list *[]string) (err error) {
	filepath.WalkDir(jc.path, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// @todo pass const Ext below instead of ".json" string
		var filename string = info.Name()
		if filepath.Ext(filename) == ".json" {
			*list = append(*list, filename[0:len(filename)-len(".json")])
		}
		return nil
	})

	return
}
