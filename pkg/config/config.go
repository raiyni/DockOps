package config

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

var Store config.Config

func init() {
	Store = *config.NewEmpty("main")
	Store.WithOptions(config.ParseEnv)
	Store.AddDriver(yaml.Driver)

	err := Store.LoadFiles("config.yml")
	if err != nil {
		panic(err.Error())
	}
}
