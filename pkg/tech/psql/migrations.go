package psql

import (
	"context"
	"github.com/Cobalt0s/creme-brulee/pkg/stateful/config"
	"github.com/golang-migrate/migrate"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"

	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

func BumpDatabaseVersion(ctx context.Context, filePath string, cfg *config.PsqlConfig) {
	log := ctxlogrus.Extract(ctx)
	// Initiate database migration
	m, err := migrate.New(filePath, cfg.GetDataSourcePSQL().String())
	if err != nil {
		log.Fatalf("failed to start db migration %v", err)
	}
	log.Info("starting db migration")
	logDbVersion(ctx, m, "before migration")
	err = m.Up()
	if err != nil {
		log.Infof("no migrations to run %v", err)
	}
	logDbVersion(ctx, m, "after migration")
	log.Info("finished db migration")
}

func logDbVersion(ctx context.Context, m *migrate.Migrate, text string) {
	log := ctxlogrus.Extract(ctx)
	version, dirty, err := m.Version()
	if err != nil {
		log.Warn(err)
	}
	log.WithFields(logrus.Fields{
		"version": version,
		"dirty":   dirty,
	}).Warn(text)
}
