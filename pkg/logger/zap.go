package logger

import (
	"github.com/bulatok/ozon-task/internal/ozon-task/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	zapDefaultEncoding = "json"
	zapLevelMapper     = map[string]zapcore.Level{
		"DEV":  zapcore.DebugLevel,
		"PROD": zapcore.InfoLevel,
		"":     zapcore.DebugLevel,
	}
)

func ProvideZap(conf *config.Config) (*zap.Logger, error) {
	zapConf := zap.Config{}
	if zapLevelMapper[conf.Service.LogLevel] == zapcore.DebugLevel {
		zapConf = zap.NewDevelopmentConfig()
	} else {
		zapConf = zap.NewProductionConfig()
	}

	zapConf.Encoding = zapDefaultEncoding
	zapConf.DisableStacktrace = true

	return zapConf.Build()
}
