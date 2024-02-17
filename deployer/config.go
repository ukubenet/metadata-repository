package deployer

import (
	"os"

	"github.com/ukubenet/metadata-repository/config"
	"github.com/ukubenet/metadata-repository/deployer/adapter"
	"github.com/ukubenet/metadata-repository/metadata"
)

const JSON = "json_file"

func get_path() string {
	path, _ := os.Getwd()
	return path + config.Config.Deployer.Path
}

func getAdapter(e metadata.EntityType) interface{} {
	if config.Config.Metadata.Adapter == JSON {
		return adapter.Local(get_path() + e.String() + "/")
	} else {
		panic("Undefined deployer adapter!")
	}
}
