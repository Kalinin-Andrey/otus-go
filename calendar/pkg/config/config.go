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
		HTTPListen              string `config:"httplisten"`
		GRPCListen              string `config:"grpclisten"`
		HTTPForPrometheuslisten string `yaml:"http_for_prometheus_listen" config:"http_for_prometheus_listen"`
	}								`config:"server"`
	Log			Log					`config:"log"`
	DB			DB					`config:"db"`
	Queue		Queue				`config:"queue"`
	Repository	struct {
		Type		string			`config:"type"`
	}								`config:"repository"`
	// JWT signing key. required.
	JWTSigningKey string			`yaml:"jwt_signing_key" config:"JWT_SIGNING_KEY,secret"`
	// JWT expiration in hours. Defaults to 72 hours (3 days)
	JWTExpiration int				`yaml:"jwt_expiration" config:"JWT_EXPIRATION"`
}

// Log is a config for a logger
type Log struct {
	Encoding		string
	OutputPaths		[]string		`config:"outputpaths"`
	Level			string
	InitialFields	map[string]interface{}	`config:"initialfields"`
}

// DB is a config for a DB connection
type DB struct {
	Dialect		string				`config:"dialect"`
	DSN			string				`config:"dsn"`
}

// Queue is a config for a queue
type Queue struct {
	RabbitMQ					RabbitMQ
	RabbitMQUserNotification	RabbitMQ
}

// RabbitMQ is a config for a connection to a RabbitMQ host
type RabbitMQ struct {
	ConsumerTag  string		`yaml:"consumer_tag"`
	URI          string
	ExchangeName string		`yaml:"exchange_name"`
	ExchangeType string		`yaml:"exchange_type"`
	Queue        string
	BindingKey   string		`yaml:"binding_key"`
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
	ctx := context.Background()

	loader := confita.NewLoader(
		file.NewBackend(pathToConfig),
		env.NewBackend(),
	)

	err := loader.Load(ctx, &config)
	if err != nil {
		return nil, err
	}

	err = confita.NewLoader(flags.NewBackend()).Load(ctx, &config)

	return &config, err
}