package psql

import (
	"context"
	config2 "github.com/Cobalt0s/creme-brulee/pkg/stateful/config"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGORM(ctx context.Context, baseConf *config2.BaseConfig, psqlConf *config2.PsqlConfig) *gorm.DB {
	log := ctxlogrus.Extract(ctx)

	dbLogLevel := logger.Silent
	if baseConf.LogLevel.String() == logrus.DebugLevel.String() {
		dbLogLevel = logger.Info
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: psqlConf.GetDataSourcePSQL().String(),
	}), &gorm.Config{
		Logger: logger.Default.LogMode(dbLogLevel),
	})
	if err != nil {
		log.Fatal(err)
	}

	return db
}
