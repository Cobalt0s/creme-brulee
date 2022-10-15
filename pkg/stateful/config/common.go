package config

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type BaseConfig struct {
	LogLevel    logrus.Level
	GinLogLevel string
	Env         string

	KafkaConfig
	PsqlConfig
}

func NewBaseConfig() (*BaseConfig, error) {
	envLogLevel, err := GetEnv("LOG_LEVEL")
	if err != nil {
		return nil, err
	}
	envEnv, err := GetEnv("ENV")
	if err != nil {
		return nil, err
	}

	logLevel := getLogLevel(envLogLevel)
	return &BaseConfig{
		LogLevel:    logLevel,
		GinLogLevel: getGinLogLevel(logLevel),
		Env:         envEnv,
	}, nil
}

func getLogLevel(logLevel string) logrus.Level {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return logrus.InfoLevel
	}
	return level
}

func getGinLogLevel(level logrus.Level) string {
	switch level {
	case logrus.TraceLevel:
		fallthrough
	case logrus.DebugLevel:
		return gin.DebugMode
	default:
		return gin.ReleaseMode
	}
}
