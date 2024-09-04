package entitystorage

import (
	"os"

	"github.com/ukubenet/metadata-repository/config"
	adapter "github.com/ukubenet/metadata-repository/entity/storage/adapter/json"
	"github.com/ukubenet/metadata-repository/global"
)

const JSON_FILE = "json_file"

func get_json_path() string {
	path, _ := os.Getwd()
	return path + config.Config.Deployer.Path + "/" + global.AppName + "/" + config.Config.Deployer.Entitysubpath + "/"
}

func getAdapter() interface{} {
	if config.Config.Deployer.Adapter == JSON_FILE {
		return adapter.JSON(get_json_path())
	} else {
		panic("Undefined adapter!")
	}
}
