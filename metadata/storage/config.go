package metastorage

import (
	"os"

	"github.com/ukubenet/metadata-repository/config"
	"github.com/ukubenet/metadata-repository/global"
	"github.com/ukubenet/metadata-repository/metadata"
	"github.com/ukubenet/metadata-repository/metadata/storage/adapter"
)

const JSON = "json_file"
const registerFolder = "register"

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

func get_register_json_path(registerType metadata.RegisterType) string {
	path, _ := os.Getwd()
	return path + config.Config.Metadata.Path + "/" + global.AppName + "/" + config.Config.Metadata.Metasubpath + "/" + registerFolder + "/" + registerType.String() + "/"
}

func getRegisterAdapter(registerType metadata.RegisterType) interface{} {
	if config.Config.Metadata.Adapter == JSON {
		return adapter.JSON(get_register_json_path(registerType))
	} else {
		panic("Undefined adapter!")
	}
}