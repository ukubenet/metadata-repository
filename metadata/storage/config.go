package metastorage

import (
	"os"

	"github.com/ukubenet/metadata-repository/config"
	"github.com/ukubenet/metadata-repository/metadata/storage/adapter"
)

const JSON = "json_file"

func get_json_path() string {
	path, _ := os.Getwd()
	return path + config.Config.Metadata.Path
}

func getAdapter() interface{} {
	if config.Config.Metadata.Adapter == JSON {
		return adapter.JSON(get_json_path())
	} else {
		panic("Undefined adapter!")
	}
}
