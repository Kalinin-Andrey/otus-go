package log

import (
	//"context"
	"github.com/pkg/errors"

	//"go.uber.org/zap/zapcore"
	"go.uber.org/zap"

	"github.com/Kalinin-Andrey/otus-go/calendar/internal/config"
)

// Logger interface
type Logger interface {
	// With returns a logger based off the root logger and decorates it with the given context and arguments.
	//With(ctx context.Context, args ...interface{}) Logger

	// Debug uses fmt.Sprint to construct and log a message at DEBUG level
	//Debug(args ...interface{})
	// Info uses fmt.Sprint to construct and log a message at INFO level
	//Info(args ...interface{})
	// Error uses fmt.Sprint to construct and log a message at ERROR level
	//Error(args ...interface{})

	// Debugf uses fmt.Sprintf to construct and log a message at DEBUG level
	//Debugf(format string, args ...interface{})
	// Infof uses fmt.Sprintf to construct and log a message at INFO level
	//Infof(format string, args ...interface{})
	// Errorf uses fmt.Sprintf to construct and log a message at ERROR level
	//Errorf(format string, args ...interface{})
}

type logger struct {
	*zap.SugaredLogger
}

// New creates a new logger
func New(conf config.Log) (loger Logger, err error) {
	var cfg zap.Config
	cfg, err = configToZapConfig(conf)
	if err != nil {
		return loger, errors.Wrapf(err, "Can not convert conf to zap conf;\nconf: %v", conf)
	}

	logger, err := cfg.Build()
	if err != nil {
		return loger, errors.Wrapf(err, "Can not build loger by cfg: %#v", cfg)
	}
	defer logger.Sync()

	logger.Info("logger construction succeeded")
	return logger, nil
}

func configToZapConfig(conf config.Log) (cfg zap.Config, err error) {
	cfg.OutputPaths = conf.OutputPaths

	level := zap.NewAtomicLevel()
	err = level.UnmarshalText([]byte(conf.Level))
	if err != nil {
		return cfg, errors.Wrapf(err, "Can not unmarshal text %q, expected one of zapcore.Levels", conf.Level)
	}
	cfg.Level = level
	return cfg,nil
}

// NewByDefault creates a new logger using the default configuration.
func NewByDefault() Logger {
	l, _ := zap.NewProduction()
	return NewWithZap(l)
}

// NewWithZap creates a new logger using the preconfigured zap logger.
func NewWithZap(l *zap.Logger) Logger {
	return &logger{l.Sugar()}
}

