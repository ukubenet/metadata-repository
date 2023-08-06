package entitystorage

import "os"

var env string = "dev"

func set_env(e string) {
	env = e
}

func get_json_path() string {
	path, _ := os.Getwd()
	if env == "test" {
		return path + "/"
	} else {
		return path + "/../../repository/json/"
	}
}
