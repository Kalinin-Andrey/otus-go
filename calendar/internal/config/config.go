package config

import(
	"context"
	"flag"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	"github.com/heetch/confita/backend/flags"
)

type Configuration struct {
	Repository	RepositoryConfiguration	`config:"repository"`
}

type RepositoryConfiguration struct {
	Type	string		`config:"type"`
}

const defaultPathToConfig = "config/config.yaml"
var pathToConfig string

var config Configuration = Configuration{

}

func Get() (error, *Configuration) {
	flag.StringVar(&pathToConfig, "conf", defaultPathToConfig, "path to YAML/JSON config file")
	flag.Parse()

	loader := confita.NewLoader(
		env.NewBackend(),
		file.NewBackend(pathToConfig),
		flags.NewBackend(),
	)
	return loader.Load(context.Background(), &config), &config
}