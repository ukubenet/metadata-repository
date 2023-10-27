package deployer

import "os"

var env string = "dev"

func SetEnv(e string) {
	env = e
}

func get_path() string {
	path, _ := os.Getwd()
	if env == "test" {
		return path + "/"
	} else {
		return path + "/../../data/entity/json/"
	}
}
