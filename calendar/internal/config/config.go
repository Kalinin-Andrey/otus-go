package config

import(
	"context"
	"flag"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	"github.com/heetch/confita/backend/flags"
)

// Configuration is the struct for app configuration
type Configuration struct {
	Repository	RepositoryConfiguration	`config:"repository"`
}

// RepositoryConfiguration is configuration of repository
type RepositoryConfiguration struct {
	Type	string		`config:"type"`
}

// defaultPathToConfig is the default path to the app config
const defaultPathToConfig = "config/config.yaml"
// pathToConfig is a path to the app config
var pathToConfig string

// config is the app config
var config Configuration = Configuration{

}

// Get func return the app config
func Get() (*Configuration, error) {
	flag.StringVar(&pathToConfig, "conf", defaultPathToConfig, "path to YAML/JSON config file")
	flag.Parse()

	loader := confita.NewLoader(
		env.NewBackend(),
		file.NewBackend(pathToConfig),
		flags.NewBackend(),
	)
	err := loader.Load(context.Background(), &config)
	return &config, err
}