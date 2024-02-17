package adapter

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ukubenet/metadata-repository/entity"
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

// Read entity from a json file
func (jc *JSONCodec) Read(entityName string, identifier string, candidate *entity.CatalogEntity) (err error) {

	file, err := os.Open(jc.path + entityName + "/" + identifier + Ext)
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

// Save the entity to JSON file
func (jc *JSONCodec) Put(candidate *entity.CatalogEntity) (err error) {
	var output []byte

	if jc.indent != "" {
		output, err = json.MarshalIndent(candidate, "", jc.indent)
	} else {
		output, err = json.Marshal(candidate)
	}

	if err != nil {
		return
	}

	os.WriteFile(jc.path+candidate.EntityName+"/"+candidate.Identifier+Ext, output, 0644)

	return
}

// Delete a json file containing entity
func (jc *JSONCodec) Delete(entityName string, identifier string) (err error) {
	err = os.Remove(jc.path + entityName + "/" + identifier + Ext)
	return
}

// Show contents of entities of specific type read from json
func (jc *JSONCodec) List(entityName string, list *[]entity.CatalogEntity) (err error) {
	filepath.WalkDir(jc.path+entityName, func(path string, info fs.DirEntry, err error) error {
		return warkpath(jc, entityName, Ext, list, path, info, err)
	})

	return
}

func warkpath(jc *JSONCodec, entityName string, ext string, list *[]entity.CatalogEntity, path string, info fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	var filename string = info.Name()
	if filepath.Ext(filename) == ext {
		var identifier = filename[0 : len(filename)-len(ext)]
		entity := new(entity.CatalogEntity)
		jc.Read(entityName, identifier, entity)
		*list = append(*list, *entity)
	}

	return nil
}

// Show list of entity types
func (jc *JSONCodec) TypeList(list *[]string) (err error) {
	files, err := ioutil.ReadDir(jc.path)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			*list = append(*list, file.Name())
		}
	}

	return
}
