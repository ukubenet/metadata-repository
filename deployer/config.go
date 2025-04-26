package deployer

import (
	"os"

	"github.com/ukubenet/metadata-repository/config"
	"github.com/ukubenet/metadata-repository/deployer/adapter"
	global "github.com/ukubenet/metadata-repository/global"
	"github.com/ukubenet/metadata-repository/metadata"
)

const JSON = "json_file"

func get_path() string {
	path, _ := os.Getwd()
	return path + config.Config.Deployer.Path + "/" + global.AppName + "/" + config.Config.Deployer.Entitysubpath
}

func getAdapter(e metadata.EntityType) interface{} {
		if config.Config.Metadata.Adapter == JSON {
		return adapter.Local(get_path() + "/" + e.String() + "/")
	} else {
		panic("Undefined deployer adapter!")
	}
}

func getRegisterAdapter(e metadata.RegisterType) interface{} {
	if config.Config.Metadata.Adapter == JSON {
	return adapter.Local(get_path() + "/register/" + e.String() + "/")
} else {
	panic("Undefined deployer register adapter!")
}
}
