package metastorage

import (
	"os"

	"github.com/ukubenet/metadata-repository/config"
	"github.com/ukubenet/metadata-repository/global"
	"github.com/ukubenet/metadata-repository/metadata"
	"github.com/ukubenet/metadata-repository/metadata/storage/adapter"
)

const JSON = "json_file"

func get_json_path(entityType metadata.EntityType) string {
	path, _ := os.Getwd()
	return path + config.Config.Metadata.Path + "/" + global.AppName + "/" + config.Config.Metadata.Metasubpath + "/" + entityType.String() + "/"
}

func getAdapter(entityType metadata.EntityType) interface{} {
	if config.Config.Metadata.Adapter == JSON {
		return adapter.JSON(get_json_path(entityType))
	} else {
		panic("Undefined adapter!")
	}
}
