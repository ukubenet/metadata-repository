package deployer

import (
	"os"

	"github.com/ukubenet/metadata-repository/config"
	"github.com/ukubenet/metadata-repository/deployer/adapter"
)

const JSON = "json_file"

func get_path() string {
	path, _ := os.Getwd()
	return path + config.Config.Metadata.Path
}

func getAdapter() interface{} {
	if config.Config.Metadata.Adapter == JSON {
		return adapter.Local(get_path())
	} else {
		panic("Undefined deployer adapter!")
	}
}
