package config

import (
	"github.com/spf13/viper"
)

type (
	AppConfig struct {
		Deployer struct {
			Adapter         string `yaml:"adapter"`
			Path            string `yaml:"path"`
			Entitysubpath   string `yaml:"entitysubpath"`
			Registersubpath string `yaml:"registersubpath"`
		} `yaml:"deployer"`
		Metadata struct {
			Adapter     string `yaml:"metadata"`
			Path        string `yaml:"path"`
			Metasubpath string `yaml:"metasubpath"`
			Tmplsubpath string `yaml:"tmplsubpath"`
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
