package config

import (
	"github.com/spf13/viper"
)

type (
	AppConfig struct {
		Deployer struct {
			Adapter string `yaml:"adapter"`
			Path    string `yaml:"path"`
		} `yaml:"deployer"`
		Metadata struct {
			Adapter  string `yaml:"adapter"`
			Path     string `yaml:"path"`
			Tmplpath string `yaml:"tmplpath"`
		} `yaml:"metadata"`
	}
)

var Config AppConfig

func LoadConfig(path string, name string) {

	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	var err error = viper.ReadInConfig()
	if err != nil {
		panic("Config not found!")
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		panic("Config parsing error!")
	}
}
