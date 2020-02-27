package log

import (
	"go.uber.org/zap/zapcore"
	"log"
	//"context"
	"github.com/pkg/errors"

	//"go.uber.org/zap/zapcore"
	"go.uber.org/zap"

	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/config"
)

// ILogger interface
type ILogger interface {
	// With returns a logger based off the root logger and decorates it with the given context and arguments.
	//With(ctx context.Context, args ...interface{}) ILogger

	// Debug uses fmt.Sprint to construct and log a message at DEBUG level
	Debug(args ...interface{})
	// Info uses fmt.Sprint to construct and log a message at INFO level
	Info(args ...interface{})
	// Error uses fmt.Sprint to construct and log a message at ERROR level
	Error(args ...interface{})

	// Debugf uses fmt.Sprintf to construct and log a message at DEBUG level
	Debugf(format string, args ...interface{})
	// Infof uses fmt.Sprintf to construct and log a message at INFO level
	Infof(format string, args ...interface{})
	// Errorf uses fmt.Sprintf to construct and log a message at ERROR level
	Errorf(format string, args ...interface{})
}

// Logger struct
type Logger struct {
	*zap.SugaredLogger
}

// New creates a new logger
func New(conf config.Log) (*Logger, error) {
	cfg, err := configToZapConfig(conf)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not convert conf to zap conf;\nconf: %v", conf)
	}

	zapLogger, err := cfg.Build()
	if err != nil {
		return nil, errors.Wrapf(err, "Can not build loger by cfg: %#v", cfg)
	}

	logger := NewWithZap(zapLogger)

	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	logger.Info("Logger construction succeeded")
	return logger, nil
}

var defaultZapConfig	= zap.Config {
	Encoding:		"json",
	EncoderConfig:	zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "",
		NameKey:        "",
		CallerKey:      "",
		StacktraceKey:  "",
		LineEnding:     "",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     nil,
		EncodeDuration: nil,
		EncodeCaller:   nil,
		EncodeName:     nil,
	},
}

func configToZapConfig(conf config.Log) (cfg zap.Config, err error) {
	cfg.OutputPaths 	= conf.OutputPaths
	cfg.Encoding		= conf.Encoding
	cfg.InitialFields	= make(map[string]interface{}, len(conf.InitialFields))

	for key, val := range conf.InitialFields {
		cfg.InitialFields[key] = val
	}

	err = cfg.Level.UnmarshalText([]byte(conf.Level))
	if err != nil {
		return cfg, errors.Wrapf(err, "Can not unmarshal text %q, expected one of zapcore.Levels", conf.Level)
	}

	return cfg,nil
}

// NewByDefault creates a new logger using the default configuration.
func NewByDefault() *Logger {
	l, _ := zap.NewProduction()
	return NewWithZap(l)
}

// NewWithZap creates a new logger using the preconfigured zap logger.
func NewWithZap(l *zap.Logger) *Logger {
	return &Logger{l.Sugar()}
}

