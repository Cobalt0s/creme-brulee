package config

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"reflect"
	"testing"
)

func TestLogLevels(t *testing.T) {

	type args struct {
		levelString string
	}
	tests := []struct {
		name     string
		args     args
		gin      string
		logLevel logrus.Level
	}{
		{
			name: "Wrong level defaults to info",
			args: args{
				levelString: "weird stuff",
			},
			gin:      gin.ReleaseMode,
			logLevel: logrus.InfoLevel,
		},
		{
			name: "Debug level",
			args: args{
				levelString: "debug",
			},
			gin:      gin.DebugMode,
			logLevel: logrus.DebugLevel,
		},
		{
			name: "Error level",
			args: args{
				levelString: "error",
			},
			gin:      gin.ReleaseMode,
			logLevel: logrus.ErrorLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			givenLogLevel := getLogLevel(tt.args.levelString)

			if !reflect.DeepEqual(givenLogLevel, tt.logLevel) {
				t.Errorf("TestLogLevels(): getLogLevel\ngot= \t%v\nwant = \t%v", givenLogLevel, tt.logLevel)
			}

			givenGin := getGinLogLevel(givenLogLevel)
			if !reflect.DeepEqual(givenGin, tt.gin) {
				t.Errorf("TestLogLevels(): getGinLogLevel\ngot= \t%v\nwant = \t%v", givenGin, tt.gin)
			}
		})
	}
}
