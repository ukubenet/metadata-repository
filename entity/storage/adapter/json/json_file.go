package adapter

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ukubenet/metadata-repository/entity"
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

// Read entity from a json file
func (jc *JSONCodec) Read(entityType metadata.EntityType, entityName string, identifier string) (e entity.Entity, err error) {

	file, err := os.Open(jc.path + entityType.String() + "/" + entityName + "/" + identifier + Ext)
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
	e = entity.GetEntityInstance(entityType)
	jsonParser.Decode(e)

	return e, nil
}

// Save the entity to JSON file
func (jc *JSONCodec) Put(e entity.Entity) (err error) {
	var output []byte

	if jc.indent != "" {
		output, err = json.MarshalIndent(e, "", jc.indent)
	} else {
		output, err = json.Marshal(e)
	}

	if err != nil {
		return
	}

	os.WriteFile(jc.path+e.GetType().String()+"/"+e.GetName()+"/"+e.GetID()+Ext, output, 0644)

	return
}

// Delete a json file containing entity
func (jc *JSONCodec) Delete(entityType metadata.EntityType, entityName string, identifier string) (err error) {
	err = os.Remove(jc.path + entityType.String() + "/" + entityName + "/" + identifier + Ext)
	return
}

// Show contents of entities of specific type read from json
func (jc *JSONCodec) List(entityType metadata.EntityType, entityName string) (list []entity.Entity, err error) {
	filepath.WalkDir(jc.path+entityType.String()+"/"+entityName, func(path string, info fs.DirEntry, err error) error {
		return warkpath(jc, entityType, entityName, Ext, &list, path, info, err)
	})

	return
}

func warkpath(jc *JSONCodec, entityType metadata.EntityType, entityName string, ext string, list *[]entity.Entity, path string, info fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	var filename string = info.Name()
	if filepath.Ext(filename) == ext {
		var identifier = filename[0 : len(filename)-len(ext)]
		var e, _ = jc.Read(entityType, entityName, identifier)
		*list = append(*list, e)
	}

	return nil
}

// Show list of entity types
func (jc *JSONCodec) TypeList(entityType metadata.EntityType) (list []string, err error) {
	files, err := ioutil.ReadDir(jc.path + entityType.String() + "/")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			list = append(list, file.Name())
		}
	}

	return
}
