package adapter

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ukubenet/metadata-repository/metadata"
)

// Read register metadate from a json file
func (jc *JSONCodec) RegisterRead(registerName string, candidate *metadata.RegisterMetadata) (err error) {
	// Ensure the path ends with a slash before appending the register name
	if jc.path[len(jc.path)-1] != '/' {
		jc.path += "/"
	}

	file, err := os.Open(jc.path + registerName + Ext)
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

	return nil
}

// Save the register metadata to JSON file
func (jc *JSONCodec) RegisterPut(candidate *metadata.RegisterMetadata) (err error) {
	var output []byte

	if jc.indent != "" {
		output, err = json.MarshalIndent(candidate, "", jc.indent)
	} else {
		output, err = json.Marshal(candidate)
	}

	if err != nil {
		return
	}

	filePath := jc.path + candidate.RegisterName + Ext
	dir := filepath.Dir(filePath)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, output, 0644)
}

// Delete a json file containing register metadata
func (jc *JSONCodec) RegisterDelete(registerName string) (err error) {
	err = os.Remove(jc.path + registerName + Ext)
	return
}

// List json files containing register metadata
func (jc *JSONCodec) RegisterList(list *[]string) (err error) {
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
