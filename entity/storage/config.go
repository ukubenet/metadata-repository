package entitystorage

import (
	"os"

	"github.com/ukubenet/metadata-repository/config"
	adapter "github.com/ukubenet/metadata-repository/entity/storage/adapter/json"
	"github.com/ukubenet/metadata-repository/metadata"
)

const JSON_FILE = "json_file"

func get_json_path(entityType metadata.EntityType) string {
	path, _ := os.Getwd()
	return path + config.Config.Deployer.Path + entityType.String() + "/"
}

func getAdapter(entityType metadata.EntityType) interface{} {
	if config.Config.Deployer.Adapter == JSON_FILE {
		return adapter.JSON(get_json_path(entityType))
	} else {
		panic("Undefined adapter!")
	}
}
