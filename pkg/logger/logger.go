package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log          *zap.Logger
	isProduction bool
)

func IsProduction() bool {
	return isProduction
}

func Init(env ...string) *zap.Logger {
	appEnv := "development"
	if len(env) > 0 && env[0] != "" {
		appEnv = env[0]
	} else if val := os.Getenv("APP_ENV"); val != "" {
		appEnv = val
	}

	isProduction = strings.ToLower(appEnv) == "production"

	var config zap.Config

	if isProduction {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "time"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05.000")
	}

	var err error
	Log, err = config.Build()
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(Log)
	return Log
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
