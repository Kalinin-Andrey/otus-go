package config

import (
	"context"
	"flag"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	"github.com/heetch/confita/backend/flags"
)

// Configuration is the struct for app configuration
type Configuration struct {
	Server		struct {
		HTTPListen	string			`config:"httplisten"`
	}									`config:"server"`
	Log			Log					`config:"log"`
	DB			DB					`config:"db"`
	Repository	struct {
		Type		string				`config:"type"`
	}									`config:"repository"`
}

// Log is config for a logger
type Log struct {
	Encoding		string
	OutputPaths		[]string		`config:"outputPaths"`
	Level			string
	InitialFields	map[string]interface{}	`config:"initialFields"`
}

// DB is config for a DB connection
type DB struct {
	DSN			string				`config:"dsn"`
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
	flag.StringVar(&pathToConfig, "config", defaultPathToConfig, "path to YAML/JSON config file")
	flag.Parse()

	loader := confita.NewLoader(
		file.NewBackend(pathToConfig),
		env.NewBackend(),
		flags.NewBackend(),
	)
	err := loader.Load(context.Background(), &config)
	return &config, err
}