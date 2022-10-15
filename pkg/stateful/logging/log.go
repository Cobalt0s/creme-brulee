package logging

import (
	"context"
	"github.com/Cobalt0s/creme-brulee/pkg/stateful/config"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
)

func EnhanceContextWithLogger(ctx context.Context, cfg *config.BaseConfig) context.Context {
	logger := logrus.StandardLogger()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetReportCaller(true)
	logger.SetLevel(cfg.LogLevel)
	return ctxlogrus.ToContext(ctx, logrus.NewEntry(logger))
}
